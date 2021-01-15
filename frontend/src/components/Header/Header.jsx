import React from 'react';
import './Header.scss';

const Header = props => (
  <div className='header text-white d-flex justify-content-between align-items-center'>
    <h2>Go Chat App</h2>
    <p className='text-right m-0'>
      {props.username ? `Hello ${props.username}!` : 'Hello!'}
    </p>
  </div>
);

export default Header;
