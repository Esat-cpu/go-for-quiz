package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"
)

type Room struct {
	id      string
	players map[string]*Player
	clients map[string]*Client
	started bool
	mu      sync.Mutex
}

type LobbyPlayer struct {
	ID       string
	Nickname string
}

var ErrGameStarted = errors.New("game has already started")

func NewRoom(id string) *Room {
	return &Room{
		id:      id,
		players: make(map[string]*Player),
		clients: make(map[string]*Client),
	}
}

func (r *Room) Join(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.players[player.ID]; exists {
		return nil
	}

	if r.started {
		return ErrGameStarted
	}

	r.players[player.ID] = player
	return nil
}

func (r *Room) GetPlayer(playerID string) (Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.players[playerID]
	if !exists {
		return Player{}, false
	}

	return *p, true
}

func (r *Room) Connect(playerID string, client *Client) bool {
	r.mu.Lock()

	p, exists := r.players[playerID]

	if !exists {
		r.mu.Unlock()
		return false
	}

	oldConn, ok := r.clients[playerID]

	r.clients[playerID] = client

	p.Online = true
	p.connGen++

	if !r.started {
		r.broadcastLobbyPlayersLocked()
	}

	r.mu.Unlock()

	if ok {
		oldConn.cancel()
	}

	log.Printf("Player connected: Nickname: [%s] ID: [%s]\n", p.Nickname, p.ID)

	return true
}

func (r *Room) Disconnect(playerID string, client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.clients[playerID] != client {
		return
	}

	p, exists := r.players[playerID]
	if !exists {
		return
	}

	delete(r.clients, playerID)

	gen := p.connGen

	time.AfterFunc(2*time.Second, func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		p, exists := r.players[playerID]
		if !exists {
			return
		}

		if p.connGen != gen {
			return
		}

		log.Printf("Player disconnected: Nickname: [%s]\n", p.Nickname)

		if r.started {
			p.Online = false
		} else {
			delete(r.players, playerID)
			r.broadcastLobbyPlayersLocked()
		}
	})
}

func (r *Room) Broadcast(msg []byte) {
	r.mu.Lock()

	for _, client := range r.clients {
		select {
		case client.send <- msg:
		default:
			defer client.cancel()
		}
	}

	r.mu.Unlock()
}

func (r *Room) broadcastLocked(msg []byte) {
	for _, client := range r.clients {
		select {
		case client.send <- msg:
		default:
			client.cancel()
		}
	}
}

func (r *Room) onlineLobbyPlayersLocked() []LobbyPlayer {
	li := make([]LobbyPlayer, 0, 10)

	for _, player := range r.players {
		if player.Online {
			li = append(li, LobbyPlayer{
				ID:       player.ID,
				Nickname: player.Nickname,
			})
		}
	}

	return li
}

func (r *Room) broadcastLobbyPlayersLocked() {
	players := r.onlineLobbyPlayersLocked()

	lpm := LobbyPlayersMessage{
		Type:    "lobby_players",
		Players: players,
	}

	msg, err := json.Marshal(lpm)
	if err != nil {
		log.Println(err)
		return
	}

	r.broadcastLocked(msg)
}
