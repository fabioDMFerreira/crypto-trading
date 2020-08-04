import React, { FunctionComponent } from 'react';
import BootstrapTable from 'react-bootstrap/Table';

interface Props {
}

const Table: FunctionComponent<Props> = ({ children }) => (
  <BootstrapTable striped bordered hover>
    {children}
  </BootstrapTable>
)

export default Table
