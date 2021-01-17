import { Component } from 'react';
import { connect } from 'react-redux';
import User from '../User/User';
import './Sidebar.scss';

class Sidebar extends Component {
  render() {
    console.log(this.props);
    return (
      <div className='vh-80 text-white' id='sidebar-wrapper'>
        <div className='sidebar-heading'>
          <p className='text-center m-0'>
            {this.props.user.name
              ? `Hello ${this.props.user.name}!`
              : this.props.user.username
              ? `Hello ${this.props.user.username}!`
              : 'Hello!'}
          </p>
        </div>
        <div className='list-group list-group-flush overflow-auto h-100'>
          {this.props.users.map(user => (
            <User key={user.id} user={user} ws={this.props.ws} />
          ))}
        </div>
      </div>
    );
  }
}

const mapeStateToProps = state => {
  console.log('STATE SIDEBAR', state);
  return {
    ws: state.ws.ws,
    user: state.user.user,
    users: state.users.users
  };
};

export default connect(mapeStateToProps, {})(Sidebar);
