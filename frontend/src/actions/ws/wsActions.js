import { CONNECT } from '../types';

export const connectWs = user => {
  console.log('USER', user);
  const socket = new WebSocket('ws://localhost:8080/ws?name=' + user.name);
  console.log('Attempting Connection...');

  socket.onopen = () => {
    console.log('Successfully Connected');
  };

  socket.onclose = event => {
    console.log('Socket Closed Connection: ', event);
  };

  socket.onerror = error => {
    console.log('Socket Error: ', error);
  };
  console.log(socket);
  return { type: CONNECT, payload: socket };
};
