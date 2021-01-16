import { CONNECT } from '../../actions/types';

const INITIAL_STATE = {
  ws: null
};

const wsReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case CONNECT:
      return {
        ...state,
        ws: action.payload
      };
    default:
      return state;
  }
};

export default wsReducer;
