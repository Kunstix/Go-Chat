import axios from 'axios';
import { SET_USERNAME, ERROR, LOGIN, LOGOUT } from '../types';

const backendUrl = 'http://localhost:8080';

export const setUsername = username => {
  return {
    type: SET_USERNAME,
    payload: username
  };
};

export const login = user => async dispatch => {
  let error;
  let token;
  try {
    const result = await axios.post(backendUrl + '/api/login', user);
    console.log('DATA', result.data);
    if (result.data.status !== 'undefined' && result.data.status === 'error') {
      error = 'Login failed';
    } else {
      token = result.data;
    }
  } catch (e) {
    error = 'Login failed';
    console.log(e);
  }
  if (error) {
    dispatch({
      type: ERROR,
      payload: error
    });
  } else {
    dispatch({
      type: LOGIN,
      payload: { ...user, token }
    });
    dispatch({
      type: ERROR,
      payload: ''
    });
  }
};

export const register = user => async dispatch => {
  let error;
  try {
    const result = await axios.post(backendUrl + '/api/register', user);
    if (result.data.status !== 'undefined' && result.data.status === 'error') {
      error = 'Register failed';
    }
  } catch (e) {
    error = 'Register failed';
    console.log(e);
  }
  if (error) {
    dispatch({
      type: ERROR,
      payload: error
    });
  } else {
    dispatch({
      type: ERROR,
      payload: ''
    });
  }
};

export const logout = () => {
  return {
    type: LOGOUT
  };
};
