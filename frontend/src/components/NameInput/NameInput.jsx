import React, { Component } from 'react';
import { connect } from 'react-redux';
import { connectWs } from '../../actions/ws/wsActions';
import { setUsername, login, register } from '../../actions/user/userActions';

class NameInput extends Component {
  constructor(props) {
    super(props);
    this.state = { name: '', username: '', password: '' };
  }

  handleChange = (target, event) => {
    this.setState({ [target]: event.target.value });
  };

  handleAnonSubmit = event => {
    event.preventDefault();
    this.props.setUsername(this.state.name);
  };

  handleLogin = event => {
    event.preventDefault();
    console.log('LOGIN', event);
    this.props.login({
      username: this.state.username,
      password: this.state.password
    });
  };

  handleRegister = event => {
    event.preventDefault();
    console.log('Register', event);
    this.props.register({
      username: this.state.username,
      password: this.state.password
    });
  };

  render() {
    return (
      <div className='d-flex flex-column align-items-center text-white'>
        <form
          className='form-group border p-4 mt-4 d-flex flex-column justify-content-center align-items-center'
          onSubmit={this.handleAnonSubmit}
        >
          <label>
            Anonymous login:
            <input
              type='text'
              value={this.state.name}
              onChange={event => this.handleChange('name', event)}
              className='form-control'
              placeholder='Join the app by adding your name'
            />
          </label>
          <input
            type='submit'
            className='btn btn-primary btn-sm ml-2'
            value='Join'
          />
        </form>
        <form
          className='form-group border p-4 d-flex flex-column justify-content-center align-items-center'
          onSubmit={this.handleLogin}
        >
          <label>
            User login:
            <input
              type='text'
              value={this.state.username}
              onChange={event => this.handleChange('username', event)}
              className='form-control'
              placeholder='Join the app by adding your name'
            />
            <input
              type='password'
              value={this.state.password}
              onChange={event => this.handleChange('password', event)}
              className='form-control'
              placeholder='Join the app by adding your name'
            />
          </label>
          <input
            type='submit'
            className='btn btn-primary btn-sm ml-2'
            value='Login'
          />
          <input
            type='submit'
            className='btn btn-primary btn-sm ml-2'
            value='Register'
            onClick={event => this.handleRegister(event)}
          />
        </form>
      </div>
    );
  }
}

const mapeStateToProps = state => {
  return {
    rooms: state.rooms,
    user: state.auth.currentUser,
    ws: state.ws
  };
};

export default connect(mapeStateToProps, {
  connectWs,
  setUsername,
  login,
  register
})(NameInput);
