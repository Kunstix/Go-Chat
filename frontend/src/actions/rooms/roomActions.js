import {
  LEAVE_ROOM,
  RECEIVE_MSG,
  USER_JOINED,
  USER_LEFT,
  ROOM_JOINED
} from '../types';

export const receiveMsg = data => {
  let msg = JSON.parse(data);
  switch (msg.action) {
    case 'send-message':
      return {
        type: RECEIVE_MSG,
        payload: msg
      };
    case 'user-join':
      return {
        type: USER_JOINED,
        payload: msg.sender
      };
    case 'user-left':
      return {
        type: USER_LEFT,
        payload: msg.sender
      };
    case 'room-joined':
      return {
        type: ROOM_JOINED,
        payload: msg
      };
    default:
      break;
  }
};

export const joinRoom = (room, ws) => {
  ws.send(JSON.stringify({ action: 'join-room', message: room }));
};

export const joinPrivateRoom = (user, ws) => {
  ws.send(JSON.stringify({ action: 'join-room-private', message: user.id }));
};

export const leaveRoom = (room, ws) => {
  ws.send(JSON.stringify({ action: 'leave-room', message: room }));
  return {
    type: LEAVE_ROOM,
    payload: room
  };
};
