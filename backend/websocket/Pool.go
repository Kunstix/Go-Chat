package websocket

import (
	"github.com/google/uuid"
	"log"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[uuid.UUID]*Client
	Broadcast  chan []byte
	Rooms      map[*Room]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan []byte, 256),
		Rooms:      make(map[*Room]bool),
	}
}

func (pool *Pool) Start() {
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

func (pool *Pool) registerClient(client *Client) {
	pool.notifyClientJoined(client)
	pool.listOnlineClients(client)
	pool.Clients[client.ID] = client
}

func (pool *Pool) unregisterClient(client *Client) {
	if _, ok := pool.Clients[client.ID]; ok {
		delete(pool.Clients, client.ID)
		pool.notifyClientLeft(client)
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
		if room.Name == name {
			return room
		}
	}
	return foundRoom
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
	go room.Run()
	pool.Rooms[room] = true

	return room
}

func (pool *Pool) notifyClientJoined(c *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: c,
	}
	log.Printf("Notify client joined... %s %s\n", UserJoinedAction, c.Name)
	pool.broadcastToClients(message.encode())
}

func (pool *Pool) notifyClientLeft(c *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: c,
	}
	log.Printf("Notify client left... %s %s\n", UserLeftAction, c.Name)
	pool.broadcastToClients(message.encode())
}

func (pool *Pool) listOnlineClients(client *Client) {
	for _, c := range pool.Clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: c,
		}
		client.Send <- message.encode()
	}
}
