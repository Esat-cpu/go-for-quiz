package main

import (
	"errors"
	"sync"
)

type Room struct {
	ID      string
	Players map[string]*Player
	Started bool
	mu      sync.Mutex
}

var ErrGameStarted = errors.New("game has already started")

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
	}
}

func (r *Room) Join(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Players[player.ID]; exists {
		return nil
	}

	if r.Started {
		return ErrGameStarted
	}

	r.Players[player.ID] = player
	return nil
}

func (r *Room) GetPlayer(playerID string) (Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.Players[playerID]
	if !exists{
		return Player{}, false
	}

	return *p, true
}
