package main

import (
	"crypto/rand"
	"encoding/hex"
)

type Player struct {
	ID       string
	Nickname string
	Score    int
	Online   bool
}

func NewPlayer(nickname string) *Player {
	var b [16]byte
	rand.Read(b[:])
	id := hex.EncodeToString(b[:])

	return &Player{
		ID:       id,
		Nickname: nickname,
	}
}
