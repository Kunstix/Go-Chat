import {
  LEAVE_ROOM,
  RECEIVE_MSG,
  USER_JOINED,
  USER_LEFT,
  ROOM_JOINED
} from '../types';
import { leavingRoom } from '../../api';

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

export const leaveRoom = (roomId, ws) => {
  leavingRoom(roomId, ws);
  return {
    type: LEAVE_ROOM,
    payload: roomId
  };
};
