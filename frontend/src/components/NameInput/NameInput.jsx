import React, { Component } from 'react';
import { connect } from 'react-redux';
import { connectWs } from '../../actions/ws/wsActions';
import { setUsername } from '../../actions/user/userActions';

class NameInput extends Component {
  constructor(props) {
    super(props);
    this.state = { username: '' };
  }

  handleChange = event => {
    this.setState({ username: event.target.value });
  };

  handleSubmit = event => {
    event.preventDefault();
    this.props.setUsername(this.state.username);
  };

  render() {
    return (
      <form
        className='form-group mt-4 d-flex justify-content-center align-items-center'
        onSubmit={this.handleSubmit}
      >
        <label className='w-50'>
          <input
            type='text'
            value={this.state.username}
            onChange={this.handleChange}
            className='form-control'
            placeholder='Join the app by adding your name'
          />
        </label>
        <input
          type='submit'
          className='btn btn-primary btn-sm ml-2'
          value='Submit'
        />
      </form>
    );
  }
}

const mapeStateToProps = state => {
  console.log(`State in mapStateToProps ${state}`, state);
  return {
    rooms: state.rooms,
    user: state.user.user,
    ws: state.ws
  };
};

export default connect(mapeStateToProps, { connectWs, setUsername })(NameInput);
