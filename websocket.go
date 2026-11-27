package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type Client struct {
	conn   *websocket.Conn
	cancel context.CancelFunc
	send   chan []byte
}

func (c *Client) writePump(ctx context.Context) {
	for {
		select {
		case msg := <-c.send:
			writeCtx, writeCancel := context.WithTimeout(ctx, 2*time.Second)

			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			writeCancel()
			if err != nil {
				c.cancel()
				log.Println(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	p, ok := a.GetPlayerFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &Client{
		conn:   conn,
		cancel: cancel,
		send:   make(chan []byte, 32),
	}

	ok = a.Room.Connect(p.ID, client)
	if !ok {
		conn.Close(websocket.StatusPolicyViolation, "Connection rejected")
		return
	}
	defer a.Room.Disconnect(p.ID, client)

	go heartbeat(ctx, conn, cancel)
	go client.writePump(ctx)

	for {
		msgType, msg, err := conn.Read(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%s %s\n", msgType, msg)
	}
}

func heartbeat(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)

			err := conn.Ping(pingCtx)
			pingCancel()

			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
