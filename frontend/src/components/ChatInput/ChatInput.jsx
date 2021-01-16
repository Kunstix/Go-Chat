import React, { Component } from 'react';
import { connect } from 'react-redux';
import { sendMsg } from '../../api';

class ChatInput extends Component {
  constructor(props) {
    super(props);
    this.state = { msg: '' };
  }

  handleChange = event => {
    event.preventDefault();
    this.setState({ msg: event.target.value });
  };

  handleSubmit = event => {
    event.preventDefault();
    sendMsg(this.state.msg, this.props.ws, this.props.room);
    this.setState({ msg: '' });
  };

  render() {
    return (
      <form className='form-group my-2' onSubmit={this.handleSubmit}>
        <label>
          <input
            type='text'
            value={this.state.msg}
            onChange={this.handleChange}
            className='form-control'
          />
        </label>
        <input
          type='submit'
          className='btn btn-primary btn-sm  ml-2'
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

export default connect(mapeStateToProps, {})(ChatInput);
