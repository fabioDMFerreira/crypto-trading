import { render } from '@testing-library/react';
import React from 'react';

import Applications from './Applications';

describe('Applications', () => {
  it('should render', () => {
    render(<Applications />);
  });
});
