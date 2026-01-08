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
	// Change this target if needed (common ones: "youtube.com:443", "www.google.com:443", "play.googleapis.com:443")
	target := "youtube.com:443"

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	// Bidirectional copy between websocket and TCP connection
	go io.Copy(conn, ws.UnderlyingConn())
	go io.Copy(ws.UnderlyingConn(), conn)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Optional: Restrict to specific path if your config expects it (e.g., "/nConnection" or "/app2")
	// if r.URL.Path != "/nConnection" { http.NotFound(w, r); return }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	proxy(ws)
}

func main() {
	// Handle root path (change "/" to "/nConnection" if your old config requires a specific path)
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
