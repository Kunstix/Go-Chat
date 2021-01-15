import { combineReducers } from 'redux';
import userReducer from './user/userReducer';
import roomsReducer from './rooms/roomsReducer';
import wsReducer from './ws/wsReducer';
import usersReducer from './user/usersReducer';

const rootReducer = combineReducers({
  ws: wsReducer,
  user: userReducer,
  users: usersReducer,
  rooms: roomsReducer
});

export default rootReducer;
