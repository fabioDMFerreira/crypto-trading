import { render } from '@testing-library/react';
import React from 'react';
import {
  BrowserRouter as Router,
} from 'react-router-dom';

import Navbar from '.';

describe('Navbar', () => {
  it('should render', () => {
    render(
      <Router>
        <Navbar />
      </Router>,
    );
  });
});
