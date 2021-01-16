import { Component } from 'react';
import { joinPrivateRoom } from '../../api';
import './User.scss';

class User extends Component {
  sendPrivateJoin = () => {
    joinPrivateRoom(this.props.user, this.props.ws);
  };

  render() {
    return (
      <button
        onClick={() => this.sendPrivateJoin()}
        className='list-group-item list-group-item-action text-white'
      >
        {this.props.user.name}
      </button>
    );
  }
}

export default User;
