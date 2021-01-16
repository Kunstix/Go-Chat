import React, { Component } from 'react';
import { connect } from 'react-redux';
import { joinRoom } from '../../api';

class RoomInput extends Component {
  constructor(props) {
    super(props);
    this.state = { room: '' };
  }

  handleChange = event => {
    this.setState({ room: event.target.value });
  };

  handleSubmit = event => {
    event.preventDefault();
    joinRoom(this.state.room, this.props.ws);
  };

  render() {
    return (
      <form className='form-group mt-4 mb-2' onSubmit={this.handleSubmit}>
        <label>
          <input
            className='form-control'
            type='text'
            value={this.state.room}
            onChange={this.handleChange}
            placeholder='Join Room'
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
  return {
    ws: state.ws.ws
  };
};

export default connect(mapeStateToProps, {})(RoomInput);
