package main

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
