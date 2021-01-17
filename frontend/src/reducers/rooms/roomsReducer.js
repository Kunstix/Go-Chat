import {
  JOIN_ROOM,
  LEAVE_ROOM,
  RECEIVE_MSG,
  ROOM_JOINED
} from '../../actions/types';

const INITIAL_STATE = {
  rooms: []
};

const roomsReducer = (state = INITIAL_STATE, action) => {
  const msg = action.payload;
  switch (action.type) {
    case JOIN_ROOM:
      return {
        ...state,
        rooms: [...state.rooms, { ...action.payload, messages: [] }]
      };
    case LEAVE_ROOM:
      console.log('LEAVE', state.rooms, action.payload);
      return {
        ...state,
        rooms: state.rooms.filter(room => room.id !== action.payload)
      };
    case RECEIVE_MSG:
      const foundRoom = findRoom(state, msg.target, msg);
      if (typeof foundRoom !== 'undefined') {
        foundRoom.messages.push(msg);
      }
      return {
        ...state,
        rooms: [...state.rooms]
      };
    case ROOM_JOINED:
      const room = msg.target;
      room.name = room.private ? msg.sender.name : room.name;
      room['messages'] = [];
      return {
        ...state,
        rooms: [...state.rooms, room]
      };
    default:
      return state;
  }
};

const findRoom = (state, room) => {
  return state.rooms.find(currentRoom => currentRoom.id === room.id);
};

export default roomsReducer;
