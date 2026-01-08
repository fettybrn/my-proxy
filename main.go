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
	// Restrict to exact path "/app62" to match your old payload GET /app62
	if r.URL.Path != "/app62" {
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
	// Only handle the specific upgrade path from your payload
	http.HandleFunc("/app62", handler)

	// Optional: Handle root with 404 to avoid wrong requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
