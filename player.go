package main

type Player struct {
	ID       string
	Nickname string
	Score    int
	Online   bool
	connGen  uint64
}

func NewPlayer(nickname string) *Player {
	return &Player{
		ID:       randHex(),
		Nickname: nickname,
	}
}
