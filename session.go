package main

import "net/http"

type Session struct {
	PlayerID string
}

func (a *App) CreateSession(playerID string) string {
	s := Session{
		PlayerID: playerID,
	}

	id := randHex()

	a.sessionMu.Lock()
	defer a.sessionMu.Unlock()

	a.Sessions[id] = s
	return id
}

func (a *App) GetSession(sessionID string) (Session, bool) {
	a.sessionMu.Lock()
	defer a.sessionMu.Unlock()

	s, exists := a.Sessions[sessionID]
	return s, exists
}

func (a *App) GetSessionFromRequest(r *http.Request) (Session, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return Session{}, false
	}

	return a.GetSession(cookie.Value)
}

func (a *App) GetPlayerFromRequest(r *http.Request) (Player, bool) {
	s, ok := a.GetSessionFromRequest(r)
	if !ok {
		return Player{}, false
	}

	return a.Room.GetPlayer(s.PlayerID)
}
