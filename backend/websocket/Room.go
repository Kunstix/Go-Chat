package websocket

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Private    bool      `json:"private"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		Private:    private,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (room *Room) Run() {
	for {
		select {
		case c := <-room.register:
			room.registerClientInRoom(c)
			log.Printf("Room is registered %s...", room.Name)
		case c := <-room.unregister:
			room.unregisterClientInRoom(c)
			log.Printf("Room is unregistered %s...", room.Name)
		case msg := <-room.broadcast:
			room.broadcastToClientsInRoom(msg.encode())
		}
	}
}

func (room *Room) notifyClientJoined(c *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, c.Name),
	}

	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) notifyClientLeft(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(goodbyeMessage, client.Name),
	}
	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) registerClientInRoom(c *Client) {
	room.notifyClientJoined(c)
	room.clients[c] = true
}

func (room *Room) unregisterClientInRoom(c *Client) {
	if _, ok := room.clients[c]; ok {
		delete(room.clients, c)
	}
	room.notifyClientLeft(c)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.Send <- message
	}
}
