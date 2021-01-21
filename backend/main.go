package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/kunstix/gochat/repository"
	"github.com/tkanos/gonfig"
	"net/http"

	"github.com/kunstix/gochat/api"
	"github.com/kunstix/gochat/auth"
	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/hub"

	"github.com/gorilla/mux"
)

func setupRoutes(db *sql.DB) *mux.Router {
	userRepository := &repository.UserRepository{Db: db}
	chatHub := hub.NewHub(&repository.RoomRepository{Db: db}, userRepository)
	go chatHub.Start()
	api := &api.API{UserRepository: userRepository}

	r := mux.NewRouter()
	r.Use(CORS)
	r.HandleFunc("/api/login", api.HandleLogin)
	r.HandleFunc("/api/register", api.HandleRegister)
	r.HandleFunc("/ws", auth.Secure(func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(chatHub, w, r)
	}))
	return r
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Origins")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		fmt.Printf("Options %s\n", r.Method)
		if r.Method == "OPTIONS" {
			fmt.Printf("OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func main() {
	var mode = flag.String("m", "dev", "choose dev or prod mode")
	flag.Parse()

	conf := config.Configuration{}
	if *mode == "prod" {
		err := gonfig.GetConf("config/production.json", &conf)
		if err != nil {
			panic(err)
		}
	} else {
		err := gonfig.GetConf("config/development.json", &conf)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Gopher Chat App")
	db := config.InitDB(conf)
	config.CreateRedisClient(conf)
	defer db.Close()
	r := setupRoutes(db)

	http.ListenAndServe(":"+conf.PORT, (r))
}
