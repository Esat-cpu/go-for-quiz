package main

type Player struct {
	ID       string
	Nickname string
	Score    int
	Online   bool
}

func NewPlayer(nickname string) *Player {
	return &Player{
		ID:       randHex(),
		Nickname: nickname,
	}
}
