import { LOGIN, LOGOUT, SET_USERNAME } from '../../actions/types';

const INITIAL_STATE = {
  currentUser: {
    name: '',
    username: '',
    password: '',
    token: ''
  }
};

const userReducer = (state = INITIAL_STATE, action) => {
  console.log('Login', action);
  switch (action.type) {
    case SET_USERNAME:
      return {
        ...state,
        currentUser: { name: action.payload }
      };
    case LOGIN:
      return {
        ...state,
        currentUser: action.payload
      };
    case LOGOUT:
      return {
        ...state,
        currentUser: {
          name: '',
          username: '',
          password: '',
          token: ''
        }
      };
    /*     case REHYDRATE:
      console.log('REHYDRATE', action);
      return {
        ...state,
        user: action.payload.user
      }; */
    default:
      return state;
  }
};

export default userReducer;
