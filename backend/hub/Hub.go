package hub

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/models"
	"log"
)

type Hub struct {
	Register       chan *Client
	Unregister     chan *Client
	Clients        map[string]*Client
	Users          []models.User
	Broadcast      chan []byte
	Rooms          map[*Room]bool
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewHub(roomRepository models.RoomRepository, userRepository models.UserRepository) *Hub {
	hub := &Hub{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make(map[string]*Client),
		Broadcast:      make(chan []byte, 256),
		Rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}
	hub.Users = userRepository.GetAllUsers()
	return hub
}

func (hub *Hub) Start() {
	go hub.listenGeneralChannel()

	for {
		select {
		case c := <-hub.Register:
			log.Printf("Hub received Register %s\n", c.Name)
			hub.registerClient(c)
		case c := <-hub.Unregister:
			log.Printf("Hub received Unregister %s\n", c.Name)
			hub.unregisterClient(c)
		case message := <-hub.Broadcast:
			log.Printf("Hub received Message\n")
			hub.broadcastToClients(message)
		}
	}
}

func (hub *Hub) listenGeneralChannel() {
	ch := config.Redis.Subscribe(ctx, GeneralChannel).Channel()
	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}
		switch message.Action {
		case UserJoinedAction:
			hub.handleUserJoined(message)
		case UserLeftAction:
			hub.handleUserLeft(message)
		case JoinRoomPrivateAction:
			hub.handleUserJoinPrivate(message)
		}
	}
}

func (hub *Hub) handleUserJoined(msg Message) {
	hub.Users = append(hub.Users, msg.Sender)
	hub.broadcastToClients(msg.encode())
}

func (hub *Hub) handleUserLeft(msg Message) {
	for i, user := range hub.Users {
		if user.GetId() == msg.Sender.GetId() {
			hub.Users[i] = hub.Users[len(hub.Users)-1]
			hub.Users = hub.Users[:len(hub.Users)-1]
			break
		}
	}
	hub.broadcastToClients(msg.encode())
}

func (hub *Hub) handleUserJoinPrivate(msg Message) {
	targetClients := hub.findClientsByID(msg.Message)
	for _, targetClient := range targetClients {
		targetClient.joinRoom(msg.Target.GetName(), msg.Sender)
	}
}

func (hub *Hub) publishClientJoined(c *Client) {
	msg := &Message{
		Action: UserJoinedAction,
		Sender: c,
	}
	if err := config.Redis.Publish(ctx, GeneralChannel, msg.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (hub *Hub) publishClientLeft(c *Client) {
	msg := &Message{
		Action: UserLeftAction,
		Sender: c,
	}
	if err := config.Redis.Publish(ctx, GeneralChannel, msg.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (hub *Hub) registerClient(c *Client) {
	if user := hub.findUserByID(c.GetId()); user == nil {
		hub.userRepository.AddUser(c)
	}

	hub.publishClientJoined(c)
	hub.listOnlineUsers(c)
	hub.Clients[c.GetId()] = c
}

func (hub *Hub) unregisterClient(c *Client) {
	if _, ok := hub.Clients[c.GetId()]; ok {
		delete(hub.Clients, c.GetId())
		hub.publishClientLeft(c)
	}
}

func (hub *Hub) broadcastToClients(message []byte) {
	log.Println("Broadcasting to clients...")
	for _, c := range hub.Clients {
		c.Send <- message
	}
}

func (hub *Hub) runRoomFromRepository(name string) *Room {
	var room *Room
	dbRoom := hub.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.Run()
		hub.Rooms[room] = true
	}

	return room
}

func (hub *Hub) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range hub.Rooms {
		if room.GetName() == name {
			return room
		}
	}

	if foundRoom == nil {
		foundRoom = hub.runRoomFromRepository(name)
	}
	return foundRoom
}

func (hub *Hub) findUserByID(ID string) models.User {
	for _, user := range hub.Users {
		if user.GetId() == ID {
			return user
		}
	}
	return nil
}

func (hub *Hub) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range hub.Rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

func (hub *Hub) findClientByID(ID string) *Client {
	var foundClient *Client
	for _, c := range hub.Clients {
		if c.GetId() == ID {
			foundClient = c
			break
		}
	}
	return foundClient
}

func (hub *Hub) findClientsByID(ID string) []*Client {
	var foundClients []*Client
	for _, c := range hub.Clients {
		if c.GetId() == ID {
			foundClients = append(foundClients, c)
		}
	}

	return foundClients
}

func (hub *Hub) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	hub.roomRepository.AddRoom(room)

	go room.Run()
	hub.Rooms[room] = true

	return room
}

func (hub *Hub) listOnlineUsers(c *Client) {
	var users = make(map[string]bool)
	for _, user := range hub.Users {
		if ok := users[user.GetId()]; !ok {
			message := &Message{
				Action: UserJoinedAction,
				Sender: user,
			}
			users[user.GetId()] = true
			c.Send <- message.encode()
		}
	}
}
