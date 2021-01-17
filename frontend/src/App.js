import React, { Component } from 'react';
import './App.css';
import NameInput from './components/NameInput/NameInput';
import RoomOverview from './components/RoomOverview/RoomOverview';
import { connect } from 'react-redux';
import Sidebar from './components/Sidebar/Sidebar';

class App extends Component {
  render() {
    return (
      <div className='d-flex App' id='wrapper'>
        <Sidebar />
        <div id='page-content-wrapper'>
          <nav className='navbar navbar-expand-lg'>
            {/*             <button className='btn btn-primary' id='menu-toggle'>
              <span className='navbar-toggler-icon'></span>
            </button> */}
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
  return {
    username: state.user.user.name,
    token: state.user.user.token
  };
};

export default connect(mapStateToProps, {})(App);
