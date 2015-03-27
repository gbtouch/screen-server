package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	namespace string
	id        int
	// Buffered channel of outbound messages.
	send chan []byte
}

func waHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	rec := map[string]interface{}{}
	for {
		if err = ws.ReadJSON(&rec); err != nil {
			if err.Error() == "EOF" {
				return
			}
			// ErrShortWrite means that a write accepted fewer bytes than requested but failed to return an explicit error.
			if err.Error() == "unexpected EOF" {
				return
			}
			fmt.Println("Read : " + err.Error())
			return
		}
		rec["Test"] = "server:i'm tommy"
		fmt.Println(rec)
		if err = ws.WriteJSON(&rec); err != nil {
			fmt.Println("Write : " + err.Error())
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", waHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
