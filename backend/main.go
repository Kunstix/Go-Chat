package main

import (
	"database/sql"
	"fmt"

	"github.com/kunstix/gochat/repository"
	"net/http"

	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/websocket"
)

func setupRoutes(db *sql.DB) {
	pool := websocket.NewPool(&repository.RoomRepository{Db: db}, &repository.UserRepository{Db: db})
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})
}

func main() {
	fmt.Println("Gopher Chat App")
	db := config.InitDB()
	config.CreateRedisClient()

	defer db.Close()
	setupRoutes(db)
	http.ListenAndServe(":8080", nil)
}
