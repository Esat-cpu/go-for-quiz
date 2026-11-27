package main

type LobbyPlayersMessage struct {
	Type    string        `json:"type"`
	Players []LobbyPlayer `json:"players"`
}
