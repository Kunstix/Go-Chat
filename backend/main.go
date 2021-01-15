package main

import (
	"fmt"
	"net/http"

	"github.com/kunstix/gochat/websocket"
)

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Chat App")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
