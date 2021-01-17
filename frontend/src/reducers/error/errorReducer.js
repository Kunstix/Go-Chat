import { ERROR } from '../../actions/types';

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
    default:
      return state;
  }
};

export default errorReducer;
