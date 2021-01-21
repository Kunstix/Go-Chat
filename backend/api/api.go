package api

import (
	"encoding/json"
	"github.com/kunstix/gochat/auth"
	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/repository"
	"log"
	"net/http"
)

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type API struct {
	UserRepository *repository.UserRepository
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser := api.UserRepository.FindUserByUsername(user.Username)
	if dbUser == nil {
		returnErrorResponse(w, "User does not exist")
		return
	}

	ok, err := auth.ComparePassword(user.Password, dbUser.Password)

	if !ok || err != nil {
		returnErrorResponse(w, "Invalid password")
		return
	}

	token, err := auth.CreateJWTToken(dbUser)

	if err != nil {
		returnErrorResponse(w, "jwt failed")
		return
	}

	w.Write([]byte(token))
}

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		returnErrorResponse(w, "Username and password required")
		return
	}

	dbUser := api.UserRepository.FindUserByUsername(user.Username)
	if dbUser != nil {
		returnErrorResponse(w, "User already exists")
		return
	}

	password, _ := auth.GeneratePassword(user.Password)
	api.UserRepository.AddRegisteredUser(user.Username, password)

	msg := &Message{
		Action: "user-registered",
	}
	jsonMsg, err := json.Marshal(msg)

	if err := config.Redis.Publish(config.Ctx, config.GeneralChannel, jsonMsg).Err(); err != nil {
		log.Println(err)
	}

	returnOkResponse(w)
}

func returnErrorResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "error","msg": "` + msg + `"}`))
}

func returnOkResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\": \"ok\"}"))
}
