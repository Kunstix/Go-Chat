import { CONNECT } from '../types';

export const connectWs = user => {
  var socket;
  if (user.token) {
    socket = new WebSocket('ws://localhost:8080/ws?bearer=' + user.token);
  } else {
    socket = new WebSocket('ws://localhost:8080/ws?name=' + user.name);
  }

  socket.onopen = () => {
    console.log('Successfully Connected');
  };

  socket.onclose = event => {
    console.log('Socket Closed Connection: ', event);
  };

  socket.onerror = error => {
    console.log('Socket Error: ', error);
  };
  return { type: CONNECT, payload: socket };
};
