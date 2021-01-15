package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID    uuid.UUID       `json:"id"`
	Name  string          `json:"name"`
	Conn  *websocket.Conn `json:"-"`
	Pool  *Pool           `json:"-"`
	Send  chan []byte     `json:"-"`
	Rooms map[*Room]bool  `json:"-"`
}

func NewClient(conn *websocket.Conn, pool *Pool, name string) *Client {
	log.Printf("New client %s\n", name)
	return &Client{
		ID:    uuid.New(),
		Name:  name,
		Conn:  conn,
		Pool:  pool,
		Send:  make(chan []byte, 256),
		Rooms: make(map[*Room]bool),
	}

}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket of client Called")
	conn, err := Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	c := NewClient(conn, pool, name[0])

	go c.Write()
	go c.Read()

	log.Printf("Registering user in pool %s\n", name[0])
	pool.Register <- c
}

func (c *Client) Read() {
	log.Printf("Starting read for %s\n", c.Name)
	defer func() {
		c.disconnect()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	for {
		_, byteMsg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		c.handleMsg(byteMsg)
	}
}

func (c *Client) Write() {
	log.Printf("Starting write for %s\n", c.Name)
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// The Pool closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			log.Println("Sending msg...")
			w.Write(message)
			logIt(message)

			// Attach queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				log.Println("Sending attached msgs...")
				w.Write(Newline)
				another := <-c.Send
				w.Write(another)
				logIt(another)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("Ping client %s\n", c.Name)
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMsg(jsonMsg []byte) {
	var msg Message
	if err := json.Unmarshal(jsonMsg, &msg); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}
	log.Printf("Incoming Message %s with %s", msg.Action, msg.Message)

	msg.Sender = c

	switch msg.Action {
	case SendMessageAction:
		roomID := msg.Target.ID.String()
		if room := c.Pool.findRoomByID(roomID); room != nil {
			room.broadcast <- &msg
		}
	case JoinRoomAction:
		c.handleJoinRoomMessage(msg)
	case LeaveRoomAction:
		c.handleLeaveRoomMessage(msg)
	case JoinRoomPrivateAction:
		c.handleJoinRoomPrivateMessage(msg)
	}
}

func (c *Client) handleJoinRoomMessage(msg Message) {
	roomName := msg.Message
	c.joinRoom(roomName, nil)
}

func (c *Client) handleJoinRoomPrivateMessage(msg Message) {
	target := c.Pool.findClientByID(msg.Message)
	if target == nil {
		return
	}

	roomName := msg.Message + c.ID.String()

	c.joinRoom(roomName, target)
	target.joinRoom(roomName, c)
}

func (c *Client) handleLeaveRoomMessage(message Message) {
	log.Printf("%s is leaving room %s... ", c.Name, message.Message)
	room := c.Pool.findRoomByName(message.Message)
	if _, ok := c.Rooms[room]; ok {
		delete(c.Rooms, room)
	}
	room.unregister <- c
}

func (c *Client) joinRoom(roomName string, sender *Client) {

	room := c.Pool.findRoomByName(roomName)
	if room == nil {
		room = c.Pool.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return
	}

	if !c.isInRoom(room) {
		c.Rooms[room] = true
		room.register <- c
		c.notifyRoomJoined(room, sender)
	}
}

func (c *Client) isInRoom(room *Room) bool {
	if _, ok := c.Rooms[room]; ok {
		return true
	}
	return false
}

func (c *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}
	c.Send <- message.encode()
}

func (c *Client) disconnect() {
	log.Printf("Disonnect %s\n", c.Name)
	c.Pool.Unregister <- c
	for room := range c.Rooms {
		room.unregister <- c
	}
	c.Conn.Close()
	close(c.Send)
}

func logIt(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error on unmarshal JSON message %s\n", err)
	}
	log.Printf("Message with %s %s\n", msg.Action, msg.Message)
}
