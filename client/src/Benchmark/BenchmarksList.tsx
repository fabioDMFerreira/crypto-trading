import React from 'react';
import Table from 'react-bootstrap/Table';
import { useTable } from 'react-table';

interface Props {
  benchmarks: any[],
  deleteBenchmark: any,
  selectBenchmark: any,
}

export default ({ benchmarks, deleteBenchmark, selectBenchmark }: Props) => {
  const columns = React.useMemo(
    () => [
      {
        Header: '',
        accessor: '_id',
        Cell: ({ value }: any) => <button type="button" onClick={() => { selectBenchmark(value); }}>Show Results</button>,
      },
      {
        Header: 'Points',
        accessor: 'input.statisticsOptions.numberOfPointsHold',
      },
      {
        Header: 'Data Source',
        accessor: 'input.dataSourceFilePath',
      },
      {
        Header: 'Initial Amount',
        accessor: 'input.accountInitialAmount',
      },
      {
        Header: 'Final Amount',
        accessor: 'output.finalAmount',
      },
      {
        Header: 'Assets Bought',
        accessor: 'output.buys',
        Cell: ({ value }: any) => <p>{value ? value.length : ''}</p>,
      },
      {
        Header: 'Assets Sold',
        accessor: 'output.sells',
        Cell: ({ value }: any) => <p>{value ? value.length : ''}</p>,
      },
      {
        Header: 'Assets Pending',
        accessor: 'output.sellsPending',
      },
      {
        Header: 'Status',
        accessor: 'status',
      },
      {
        Header: 'Created At',
        accessor: 'createdAt',
      },
      {
        Header: 'Completed At',
        accessor: 'completedAt',
      },
      {
        id: 'edit',
        accessor: '_id',
        Cell: ({ value }: any) => (
          <button
            type="button"
            onClick={() => {
              deleteBenchmark(value);
            }}
          >
            -
          </button>
        ),
      },
    ],
    [deleteBenchmark, selectBenchmark],
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable({ columns, data: benchmarks || [] });

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
