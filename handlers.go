package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (a *App) handleJoin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nickname string `json:"nickname"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	body.Nickname = strings.TrimSpace(body.Nickname)
	
	if body.Nickname == "" {
		http.Error(w, "nickname alanı gereklidir", http.StatusBadRequest)
		return
	}

	p := NewPlayer(body.Nickname)

	token := a.CreateSession(p.ID)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	err = a.Room.Join(p)
	if err != nil {
		if errors.Is(err, ErrGameStarted) {
			http.Error(w, "Oyun zaten başladı", http.StatusConflict)
			return
		}
		http.Error(w, "Server hatası", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
