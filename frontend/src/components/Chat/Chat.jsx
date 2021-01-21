import React, { Component } from 'react';
import { connect } from 'react-redux';
import Message from '../Message/Message';
import { leaveRoom } from '../../actions/rooms/roomActions';
import './Chat.scss';

class Chat extends Component {
  leaveChat = () => {
    this.props.leaveRoom(this.props.room.id, this.props.ws);
  };

  render() {
    const messages = this.props.room.messages.map((msg, index) => {
      return (
        <Message key={index} message={msg} username={this.props.username} />
      );
    });

    return (
      <div className='h-100'>
        <h6 className='card-title p-2 text-left text-white d-flex justify-content-between mb-0 border-bottom border-white'>
          {this.props.room.name}
          <button
            type='button'
            onClick={() => this.leaveChat()}
            className='close text-white'
            aria-label='Close'
          >
            <span className='align-self-start' aria-hidden='true'>
              &times;
            </span>
          </button>
        </h6>
        <div className='messages d-flex align-items-start flex-column'>
          {messages}
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    username: state.auth.currentUser.name,
    rooms: state.rooms.rooms,
    ws: state.ws.ws
  };
};

export default connect(mapStateToProps, { leaveRoom })(Chat);
