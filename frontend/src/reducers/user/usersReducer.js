import { USER_JOINED, USER_LEFT } from '../../actions/types';

const INITIAL_STATE = {
  users: []
};

const usersReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case USER_JOINED:
      if (state.users.some(user => user.id === action.payload.id)) {
        return {
          ...state
        };
      } else {
        return {
          ...state,
          users: [...state.users, action.payload]
        };
      }

    case USER_LEFT:
      return {
        ...state,
        users: state.users.filter(user => user.id !== action.payload.id)
      };
    default:
      return state;
  }
};

export default usersReducer;
