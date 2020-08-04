import React from 'react';

import { LogEvent } from '../types';
import Table from '../components/Table';

interface Props {
  logEvents: LogEvent[]
}

export default ({ logEvents }: Props) => {
  return (
    <Table>
      <thead>
        <th>ID</th>
        <th>Event Name</th>
        <th>Message</th>
        <th>Created At</th>
      </thead>
      <tbody>
        {
          logEvents.map(
            event=>(
              <tr>
                <td>{event._id}</td>
                <td>{event.eventName}</td>
                <td>{event.message}</td>
                <td>{event.createdAt}</td>
              </tr>
            )
          )
        }
      </tbody>
    </Table>
  );
}
