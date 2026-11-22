package main

import (
	"log"
	"net/http"
	"sync"
)

type App struct {
	Room      *Room
	Sessions  map[string]Session
	sessionMu sync.Mutex
}

func HTTPServer() {
	app := &App{
		Room:     NewRoom("1"),
		Sessions: make(map[string]Session),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.handleIndex)
	mux.HandleFunc("POST /join", app.handleJoin)
	mux.HandleFunc("GET /lobby", app.handleLobby)

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
