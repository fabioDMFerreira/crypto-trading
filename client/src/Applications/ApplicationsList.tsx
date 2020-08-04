import React from 'react';
import Table from 'react-bootstrap/Table';
import { useTable } from 'react-table';

interface Props {
  applications: any[]
  selectApplication: (value: string) => void
  deleteApplication: (value: string) => void
}

export default ({ selectApplication, deleteApplication, applications }: Props) => {
  const columns = React.useMemo(
    () => [
      {
        Header: '',
        accessor: '_id',
        Cell: ({ value }: any) => <button type="button" onClick={() => { selectApplication(value); }}>Show</button>,
      },
      {
        Header: 'Asset',
        accessor: 'asset',
      },
      {
        Header: 'Created At',
        accessor: 'createdAt',
        Cell: ({ value }: any) => (<span>{value.toString()}</span>),
      },
      {
        id: 'edit',
        accessor: '_id',
        Cell: ({ value }: any) => (
          <button
            type="button"
            onClick={() => {
              deleteApplication(value);
            }}
          >
            -
          </button>
        ),
      },
    ],
    [deleteApplication, selectApplication],
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable({ columns, data: applications || [] });

  return (
    <Table striped bordered hover size="sm" {...getTableProps()}>
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th {...column.getHeaderProps()}>{column.render('Header')}</th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()}>
              {row.cells.map((cell) => <td {...cell.getCellProps()}>{cell.render('Cell')}</td>)}
            </tr>
          );
        })}
      </tbody>
    </Table>
  );
};
