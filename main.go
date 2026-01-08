package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func proxy(ws *websocket.Conn) {
	// This matches your old [Remote Proxy]: www.google.com:443
	// (the websocket tunnel forwards all traffic here)
	target := "www.google.com:443"

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	// Bidirectional piping
	go io.Copy(conn, ws.UnderlyingConn())
	go io.Copy(ws.UnderlyingConn(), conn)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/nConnection" {
		http.NotFound(w, r)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	proxy(ws)
}

func main() {
	http.HandleFunc("/nConnection", handler)

	// Optional: 404 everything else
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Rest of main unchanged...
}
