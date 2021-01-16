package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/models"
	"log"
)

type Pool struct {
	Register       chan *Client
	Unregister     chan *Client
	Clients        map[string]*Client
	Users          []models.User
	Broadcast      chan []byte
	Rooms          map[*Room]bool
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewPool(roomRepository models.RoomRepository, userRepository models.UserRepository) *Pool {
	pool := &Pool{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make(map[string]*Client),
		Broadcast:      make(chan []byte, 256),
		Rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}
	pool.Users = userRepository.GetAllUsers()
	return pool
}

func (pool *Pool) Start() {
	go pool.listenGeneralChannel()

	for {
		select {
		case c := <-pool.Register:
			log.Printf("Pool received Register %s\n", c.Name)
			pool.registerClient(c)
		case c := <-pool.Unregister:
			log.Printf("Pool received Unregister %s\n", c.Name)
			pool.unregisterClient(c)
		case message := <-pool.Broadcast:
			log.Printf("Pool received Message\n")
			pool.broadcastToClients(message)
		}
	}
}

func (pool *Pool) listenGeneralChannel() {
	ch := config.Redis.Subscribe(ctx, GeneralChannel).Channel()
	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}
		switch message.Action {
		case UserJoinedAction:
			pool.handleUserJoined(message)
		case UserLeftAction:
			pool.handleUserLeft(message)
		case JoinRoomPrivateAction:
			pool.handleUserJoinPrivate(message)
		}
	}
}

func (pool *Pool) handleUserJoined(msg Message) {
	pool.Users = append(pool.Users, msg.Sender)
	pool.broadcastToClients(msg.encode())
}

func (pool *Pool) handleUserLeft(msg Message) {
	for i, user := range pool.Users {
		if user.GetId() == msg.Sender.GetId() {
			pool.Users[i] = pool.Users[len(pool.Users)-1]
			pool.Users = pool.Users[:len(pool.Users)-1]
		}
	}
	pool.broadcastToClients(msg.encode())
}

func (pool *Pool) handleUserJoinPrivate(msg Message) {
	targetClient := pool.findClientByID(msg.Message)
	if targetClient != nil {
		targetClient.joinRoom(msg.Target.GetName(), msg.Sender)
	}
}

func (pool *Pool) publishClientJoined(c *Client) {
	msg := &Message{
		Action: UserJoinedAction,
		Sender: c,
	}
	if err := config.Redis.Publish(ctx, GeneralChannel, msg.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (pool *Pool) publishClientLeft(c *Client) {
	msg := &Message{
		Action: UserLeftAction,
		Sender: c,
	}
	if err := config.Redis.Publish(ctx, GeneralChannel, msg.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (pool *Pool) registerClient(client *Client) {
	pool.userRepository.AddUser(client)

	pool.publishClientJoined(client)
	pool.listOnlineClients(client)
	pool.Clients[client.GetId()] = client
}

func (pool *Pool) unregisterClient(c *Client) {
	if _, ok := pool.Clients[c.GetId()]; ok {
		delete(pool.Clients, c.GetId())
		pool.userRepository.RemoveUser(c)
		pool.publishClientLeft(c)
	}
}

func (pool *Pool) broadcastToClients(message []byte) {
	log.Println("Broadcasting to clients...")
	for _, c := range pool.Clients {
		c.Send <- message
	}
}

func (pool *Pool) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range pool.Rooms {
		if room.GetName() == name {
			return room
		}
	}

	if foundRoom == nil {
		foundRoom = pool.runRoomFromRepository(name)
	}

	return foundRoom
}

func (pool *Pool) findUserByID(ID string) models.User {
	for _, client := range pool.Users {
		if client.GetId() == ID {
			return client
		}
	}
	return nil
}

func (pool *Pool) runRoomFromRepository(name string) *Room {
	var room *Room
	dbRoom := pool.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.Run()
		pool.Rooms[room] = true
	}

	return room
}

func (pool *Pool) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range pool.Rooms {
		if room.ID.String() == ID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

func (pool *Pool) findClientByID(ID string) *Client {
	var foundClient *Client
	for _, c := range pool.Clients {
		if c.ID.String() == ID {
			foundClient = c
			break
		}
	}
	return foundClient
}

func (pool *Pool) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	pool.roomRepository.AddRoom(room)

	go room.Run()
	pool.Rooms[room] = true

	return room
}

func (pool *Pool) listOnlineClients(client *Client) {
	for _, user := range pool.Users {
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.Send <- message.encode()
	}
}
