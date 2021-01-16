import React, { Component } from 'react';
import { connect } from 'react-redux';
import Chat from '../Chat/Chat';
import ChatInput from '../ChatInput/ChatInput';
import RoomInput from '../RoomInput/RoomInput';
import { connectWs } from '../../actions/ws/wsActions';
import { receiveMsg } from '../../actions/rooms/roomActions';

class RoomOverview extends Component {
  componentDidMount(state) {
    this.props.connectWs(this.props.user);
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevProps.ws !== this.props.ws) {
      this.props.ws.onmessage = msg => {
        msg.data.split(/\r?\n/).forEach(data => this.props.receiveMsg(data));
      };
    }
  }

  render() {
    return (
      <div>
        <RoomInput />
        <div className='d-flex flex-wrap justify-content-around align-items-around'>
          {this.props.rooms.map(room => {
            return (
              <div
                key={room.name}
                className='card chat mb-4'
                style={{ width: 300 }}
              >
                <Chat room={room} />
                <ChatInput room={room} />
              </div>
            );
          })}
        </div>
      </div>
    );
  }
}

const mapeStateToProps = state => {
  return {
    ws: state.ws.ws,
    rooms: state.rooms.rooms,
    user: state.user.user
  };
};

export default connect(mapeStateToProps, { receiveMsg, connectWs })(
  RoomOverview
);
