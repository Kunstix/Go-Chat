import { SET_USERNAME } from '../../actions/types';

const INITIAL_STATE = {
  user: {
    name: ''
  }
};

const userReducer = (state = INITIAL_STATE, action) => {
  console.log('USERREDUCER', action);
  switch (action.type) {
    case SET_USERNAME:
      return {
        ...state,
        user: { name: action.payload }
      };
    default:
      return state;
  }
};

export default userReducer;
