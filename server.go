package main

import (
	"fmt"
	"log"
	"net/http"
)

type App struct {
	Room *Room
}

func helo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello!")
}

func HTTPServer() {
	app := &App{
		Room: NewRoom("1"),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", helo)
	mux.HandleFunc("POST /join", app.handleJoin)

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
