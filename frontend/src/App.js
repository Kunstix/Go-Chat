import React, { Component } from 'react';
import './App.css';
import NameInput from './components/NameInput/NameInput';
import RoomOverview from './components/RoomOverview/RoomOverview';
import { connect } from 'react-redux';
import Sidebar from './components/Sidebar/Sidebar';
import { logout } from './actions/user/userActions';

class App extends Component {
  render() {
    return (
      <div className='d-flex App' id='wrapper'>
        <Sidebar />
        <div id='page-content-wrapper'>
          <nav className='navbar navbar-expand-lg'>
            {/* <button className='btn btn-primary btn' id='menu-toggle'>
              <span className='navbar-toggler-icon'></span>
            </button> */}
            <button
              className='btn btn-primary btn-sm ml-auto'
              onClick={() => this.props.logout()}
            >
              <span>Logout</span>
            </button>
          </nav>
          <div className='container'>
            {this.props.username || this.props.token ? (
              <RoomOverview />
            ) : (
              <NameInput />
            )}
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  console.log('APP', state);
  const username = state.auth.currentUser ? state.auth.currentUser.name : '';
  const token = state.auth.currentUser ? state.auth.currentUser.token : '';
  return {
    username: username,
    token: token
  };
};

export default connect(mapStateToProps, { logout })(App);
