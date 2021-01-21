import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import React from 'react';
import './App.css';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import reduxThunk from 'redux-thunk';
import rootReducer from './reducers/rootReducer';
import { persistStore } from 'redux-persist';
import { PersistGate } from 'redux-persist/integration/react';

const middlewares = [reduxThunk];
const store = createStore(rootReducer, applyMiddleware(...middlewares));
const persistor = persistStore(store);

ReactDOM.render(
  <Provider store={store}>
    <PersistGate persistor={persistor} loading={<span>Loading!</span>}>
      <App />
    </PersistGate>
  </Provider>,
  document.getElementById('root')
);

reportWebVitals();
