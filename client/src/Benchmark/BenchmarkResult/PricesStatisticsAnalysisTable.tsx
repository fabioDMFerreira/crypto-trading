import React from 'react';
import Table from '../../components/Table';

interface Props {
  prices: [number, number][]
  growth: [number, number][]
  growthOfGrowth: [number, number][]
}

const getStatistics = (points: number[]) => {
  const pointsSum = points.reduce((a, b) => a + b, 0);
  const average = pointsSum / points.length;

  const stdSum = points.reduce((a, b) => a + ((b - average) ** 2), 0);
  const std = Math.sqrt(stdSum / points.length);

  return {
    average,
    std,
  };
};

export default ({ prices, growth, growthOfGrowth }: Props) => {
  const pricesStatistics = getStatistics(prices.map((el) => el[1]));
  const growthStatistics = getStatistics(growth.map((el) => el[1]));
  const growthOfGrowthStatistics = getStatistics(growthOfGrowth.map((el) => el[1]));

  return (
    <Table>
      <thead>
        <tr>
          <th>Metric</th>
          <th>Average</th>
          <th>Std</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Price</td>
          <td>{pricesStatistics.average}</td>
          <td>{pricesStatistics.std}</td>
        </tr>
        <tr>
          <td>Growth</td>
          <td>{growthStatistics.average}</td>
          <td>{growthStatistics.std}</td>
        </tr>
        <tr>
          <td>Growth of Growth</td>
          <td>{growthOfGrowthStatistics.average}</td>
          <td>{growthOfGrowthStatistics.std}</td>
        </tr>
      </tbody>
    </Table>
  );
};
