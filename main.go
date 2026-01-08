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
	// Target to tunnel to (change if your config needs a different one, e.g., "play.googleapis.com:443")
	target := "www.google.com:443"

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	go io.Copy(conn, ws.UnderlyingConn())
	go io.Copy(ws.UnderlyingConn(), conn)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Change path if your old config used something like "/nConnection" or "/app2"
	// e.g., if r.URL.Path != "/nConnection" { http.NotFound(w, r); return }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()
	proxy(ws)
}
func main() {
	http.HandleFunc("/", handler)  // Or "/nConnection" to match exactly

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	http.HandleFunc("/", handler)  // Or change to "/nConnection" if needed later

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
