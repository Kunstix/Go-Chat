package hub

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kunstix/gochat/auth"
	"github.com/kunstix/gochat/config"
	"github.com/kunstix/gochat/models"
	"github.com/kunstix/gochat/websocks"
	"log"
	"net/http"
	"time"
)

type Client struct {
	ID    uuid.UUID       `json:"id"`
	Name  string          `json:"name"`
	Conn  *websocket.Conn `json:"-"`
	Hub   *Hub            `json:"-"`
	Send  chan []byte     `json:"-"`
	Rooms map[*Room]bool  `json:"-"`
}

func NewClient(conn *websocket.Conn, hub *Hub, name string, ID string) *Client {
	log.Printf("New client %s\n", name)
	c := &Client{
		ID:    uuid.New(),
		Name:  name,
		Conn:  conn,
		Hub:   hub,
		Send:  make(chan []byte, 256),
		Rooms: make(map[*Room]bool),
	}
	if ID != "" {
		c.ID, _ = uuid.Parse(ID)
	}
	return c

}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket of client Called")

	// get user
	ctxUser := r.Context().Value(auth.UserContextKey)
	if ctxUser == nil {
		log.Println("Not authenticated")
		return
	}
	user := ctxUser.(models.User)

	// establish websocket
	conn, err := websocks.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}

	// init
	c := NewClient(conn, hub, user.GetName(), user.GetId())

	go c.Write()
	go c.Read()

	log.Printf("Registering user in hub %s\n", user.GetName(), user.GetId())
	hub.Register <- c
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
				// The Hub closed the channel.
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
		if room := c.Hub.findRoomByID(roomID); room != nil {
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
	target := c.Hub.findUserByID(msg.Message)
	if target == nil {
		return
	}

	roomName := msg.Message + c.GetId()

	privateRoom := c.joinRoom(roomName, target)

	if privateRoom != nil {
		c.inviteTarget(target, privateRoom)
	}
}

func (c *Client) handleLeaveRoomMessage(msg Message) {
	log.Printf("%s is leaving room %s... ", c.Name, msg.Message)
	room := c.Hub.findRoomByID(msg.Message)
	if room == nil {
		return
	}
	if _, ok := c.Rooms[room]; ok {
		delete(c.Rooms, room)
	}
	room.unregister <- c
}

func (c *Client) joinRoom(roomName string, sender models.User) *Room {
	room := c.Hub.findRoomByName(roomName)
	if room == nil {
		room = c.Hub.createRoom(roomName, sender != nil)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !c.isInRoom(room) {
		c.Rooms[room] = true
		room.register <- c
		c.notifyRoomJoined(room, sender)
	}

	return room
}

func (c *Client) inviteTarget(target models.User, room *Room) {
	inviteMsg := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetId(),
		Target:  room,
		Sender:  c,
	}

	if err := config.Redis.Publish(ctx, GeneralChannel, inviteMsg.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (c *Client) isInRoom(room *Room) bool {
	if _, ok := c.Rooms[room]; ok {
		return true
	}
	return false
}

func (c *Client) notifyRoomJoined(room *Room, sender models.User) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}
	c.Send <- message.encode()
}

func (c *Client) disconnect() {
	log.Printf("Disonnect %s\n", c.Name)
	c.Hub.Unregister <- c
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

func (c *Client) GetId() string {
	return c.ID.String()
}

func (c *Client) GetName() string {
	return c.Name
}
