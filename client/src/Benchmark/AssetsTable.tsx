import React from 'react';
import Table from 'react-bootstrap/Table';

import { Asset } from '../types';

interface Props {
  assets: Asset[]
}

export default ({ assets }: Props) => (
  <Table striped bordered hover>
    <thead>
      <tr>
        <th>Amount</th>
        <th>Buy Time</th>
        <th>BuyPrice</th>
        <th>Sell Time</th>
        <th>Sell Price</th>
        <th>Profit</th>
        <th>Return</th>
      </tr>
    </thead>
    <tbody>
      {
        assets.map(
          (asset) => {
            const profit = (asset.sellPrice * asset.amount) - (asset.buyPrice * asset.amount);
            const returns = (profit / (asset.buyPrice * asset.amount)) * 100;

            return (
              <tr key={Math.random()}>
                <td>{asset.amount}</td>
                <td>{asset.buyTime}</td>
                <td>{asset.buyPrice}</td>
                <td>{asset.sold ? asset.sellTime : ''}</td>
                <td>{asset.sold ? asset.sellPrice : ''}</td>
                <td>{profit.toFixed(2)}</td>
                <td>{`${(returns).toFixed(2)}%`}</td>
              </tr>
            );
          },
        )
      }
    </tbody>
  </Table>
);
