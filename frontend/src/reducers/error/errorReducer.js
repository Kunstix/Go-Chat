import { ERROR, LOGOUT } from '../../actions/types';

const INITIAL_STATE = {
  error: ''
};

const errorReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ERROR:
      return {
        ...state,
        error: { error: action.payload }
      };
    case LOGOUT:
      return INITIAL_STATE;
    default:
      return state;
  }
};

export default errorReducer;
