import { LOGIN, SET_USERNAME } from '../../actions/types';

const INITIAL_STATE = {
  user: {
    name: '',
    username: '',
    password: '',
    token: ''
  }
};

const userReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case SET_USERNAME:
      return {
        ...state,
        user: { name: action.payload }
      };
    case LOGIN:
      return {
        ...state,
        user: action.payload
      };
    default:
      return state;
  }
};

export default userReducer;
