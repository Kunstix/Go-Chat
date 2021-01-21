import { CONNECT, LOGOUT, SET_USERNAME } from '../types';

export const connectWs = (user, cb) => dispatch => {
  var socket;
  if (user.token) {
    socket = new WebSocket('ws://localhost:8080/ws?bearer=' + user.token);
  } else {
    socket = new WebSocket('ws://localhost:8080/ws?name=' + user.name);
  }

  socket.onopen = () => {
    console.log('Successfully Connected');
    dispatch({ type: CONNECT, payload: socket });
  };

  socket.onclose = event => {
    console.log('Socket Closed Connection: ', event);
    dispatch({
      type: LOGOUT
    });
  };

  socket.onerror = error => {
    console.log('Socket Error: ', error);
    dispatch({
      type: LOGOUT
    });
  };
  dispatch({ type: 'Loading', payload: '' });
};
