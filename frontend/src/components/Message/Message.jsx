import React, { Component } from 'react';
import './Message.scss';

class Message extends Component {
  render() {
    console.log('MSG', this.props.message);
    return (
      <div
        className={`Message rounded-pill p-1 px-2 mx-2 my-1 text-white text-monospace text-light small ${
          this.props.message.sender
            ? this.props.username === this.props.message.sender.name
              ? 'bg-success align-self-end'
              : 'bg-primary'
            : 'bg-secondary align-self-stretch'
        }`}
      >
        {this.props.message.sender ? (
          <span className='small text-secondary text-right d-inline-block mr-1'>
            {this.props.message.sender.name + ':'}
          </span>
        ) : (
          <span></span>
        )}
        <span className=''>{this.props.message.message}</span>
      </div>
    );
  }
}

export default Message;
