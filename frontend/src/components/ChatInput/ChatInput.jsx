import React, { Component } from 'react';
import { connect } from 'react-redux';

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
    this.sendMsg();
    console.log('Resetting...');
    this.setState({ msg: '' });
  };

  sendMsg = () => {
    console.log(
      'Send msg: ',
      JSON.stringify({
        action: 'send-message',
        message: this.state.msg,
        target: {
          id: this.props.room.id,
          name: this.props.room.name
        }
      })
    );
    if (this.state.msg !== '') {
      this.props.ws.send(
        JSON.stringify({
          action: 'send-message',
          message: this.state.msg,
          target: {
            id: this.props.room.id,
            name: this.props.room.name
          }
        })
      );
    }
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
