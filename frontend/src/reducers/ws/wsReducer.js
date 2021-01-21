import { CONNECT, LOGOUT } from '../../actions/types';

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
    case LOGOUT:
      console.log('Closing', state);
      if (state.ws) {
        state.ws.close();
        console.log('Closed');
      }
      return INITIAL_STATE;
    default:
      return state;
  }
};

export default wsReducer;
