import axios from 'axios';
import { SET_USERNAME, ERROR, LOGIN } from '../types';

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
    const result = await axios.post('http://localhost:8080/api/login', user);
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
  }
};
