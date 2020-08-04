import React from 'react';

interface Props {
  json: any
}

export default ({ json }: Props) => (
  <code>
    {JSON.stringify(json, undefined, 2)}
  </code>
);
