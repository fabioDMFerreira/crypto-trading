import React from 'react';
import Table from 'react-bootstrap/Table';

import formatDateTime from '../../formatters/formatDateTime';

interface Props {
  prices: [number, number][]
  growth: [number, number][]
  growthOfGrowth: [number, number][]
}

interface PriceRow {
  date: number,
  price: number,
  growth: number,
  growthOfGrowth: number
}


function MapToPriceRow(prices: [number, number][], growth: [number, number][], growthOfGrowth: [number, number][]) {
  const mapper: any = {};

  prices.forEach((price) => {
    mapper[price[0]] = { price: price[1] };
  });

  growth.forEach((growth) => {
    const date = growth[0];

    if (date in mapper) {
      mapper[date].growth = growth[1];
    }
  });

  growthOfGrowth.forEach((growthOfGrowth) => {
    const date = growthOfGrowth[0];

    if (date in mapper) {
      mapper[date].growthOfGrowth = growthOfGrowth[1];
    }
  });

  return Object.keys(mapper).map((date) => ({
    date,
    ...mapper[date],
  }), []);
}

export default ({ prices, growth, growthOfGrowth }: Props) => {
  const data: PriceRow[] = MapToPriceRow(prices, growth, growthOfGrowth);

  return (
    <Table striped bordered hover>
      <thead>
        <tr>
          <th>Date</th>
          <th>Price</th>
          <th>Growth</th>
          <th>Growth of Growth</th>
        </tr>
      </thead>
      <tbody>
        {
          data.map(
            (pr) => (
              <tr key={Math.random()}>
                <td>{formatDateTime(pr.date)}</td>
                <td>{pr.price}</td>
                <td style={{ backgroundColor: pr.growth > 0 ? '#C8E6C9' : '#FFCDD2' }}>{pr.growth}</td>
                <td style={{ backgroundColor: pr.growthOfGrowth > 0 ? '#C8E6C9' : '#FFCDD2' }}>{pr.growthOfGrowth}</td>
              </tr>
            ),
          )
        }
      </tbody>
    </Table>
  );
};
