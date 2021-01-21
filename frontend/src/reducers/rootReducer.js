import { combineReducers } from 'redux';
import userReducer from './user/userReducer';
import usersReducer from './user/usersReducer';
import roomsReducer from './rooms/roomsReducer';
import wsReducer from './ws/wsReducer';
import errorReducer from './error/errorReducer';
import { persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import { LOGOUT } from '../actions/types';
import autoMergeLevel2 from 'redux-persist/lib/stateReconciler/autoMergeLevel2';

const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['auth', 'rooms', 'users', 'ws'],
  stateReconciler: autoMergeLevel2
};

const rootReducer = (state, action) => {
  /*   if (action.type === LOGOUT) {
    localStorage.removeItem('persist:root');
  } */
  return appReducer(state, action);
};

const appReducer = combineReducers({
  ws: wsReducer,
  auth: userReducer,
  users: usersReducer,
  rooms: roomsReducer,
  error: errorReducer
});

export default persistReducer(persistConfig, rootReducer);
