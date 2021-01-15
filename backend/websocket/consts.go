package websocket

import (
	"time"
)

const (
	// Max wait time when writing message to peer
	WriteWait = 10 * time.Second

	// Max time till next pong from peer
	PongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 10000
)

const (
	SendMessageAction     = "send-message"
	JoinRoomAction        = "join-room" // join public room
	LeaveRoomAction       = "leave-room"
	UserJoinedAction      = "user-join"
	UserLeftAction        = "user-left"
	JoinRoomPrivateAction = "join-room-private" // create private room
	RoomJoinedAction      = "room-joined"       // answer to create room
)
const welcomeMessage = "%s joined the room"
const goodbyeMessage = "%s left the room"

var (
	Newline = []byte{'\n'}
	Space   = []byte{' '}
)
