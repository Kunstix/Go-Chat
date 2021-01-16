export const joinRoom = (room, ws) => {
  ws.send(JSON.stringify({ action: 'join-room', message: room }));
};

export const joinPrivateRoom = (user, ws) => {
  ws.send(JSON.stringify({ action: 'join-room-private', message: user.id }));
};

export const leavingRoom = (room, ws) => {
  ws.send(JSON.stringify({ action: 'leave-room', message: room }));
};

export const sendMsg = (msg, ws, room) => {
  if (msg !== '') {
    ws.send(
      JSON.stringify({
        action: 'send-message',
        message: msg,
        target: {
          id: room.id,
          name: room.name
        }
      })
    );
  }
};
