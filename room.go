package main

import (
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
	r.mu.Unlock()

	if ok {
		oldConn.cancel()
	}

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
		}
	})
}
