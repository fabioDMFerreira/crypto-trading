import { render } from '@testing-library/react';
import React from 'react';

import { Application } from '../types';
import ApplicationsList from './ApplicationsList';

const applications: Application[] = [
  {
    _id: '5f280d4c306a86d704a77fdb',
    asset: 'ETH',
    accountID: '5f280d4c306a86d704a77fda',
    options: {
      notificationOptions: {
        receiver: 'martinhoferreira10@gmail.com',
        sender: 'devlogin18@gmail.com',
        senderpassword: 'dlifgkgwatyidrcz',
      },
      statisticsOptions: {
        numberOfPointsHold: 5000,
      },
      decisionMakerOptions: {
        maximumBuyAmount: 0.0,
        maximumFIATBuyAmount: 500.0,
        minimumProfitPerSold: 0.00999999977648258,
        minimumPriceDropToBuy: 0.00999999977648258,
        growthDecreaseLimit: -100.0,
        growthIncreaseLimit: 100.0,
      },
      collectorOptions: {
        priceVariationDetection: 0.00999999977648258,
        datasource: null,
        newPriceTimeRate: 15,
      },
    },
    createdAt: new Date(2020, 1, 1),
  },
  {
    _id: '5f280a2b8ab05b01afea4703',
    asset: 'BTC',
    accountID: '5f280a2b8ab05b01afea4702',
    options: {
      notificationOptions: {
        receiver: 'martinhoferreira10@gmail.com',
        sender: 'devlogin18@gmail.com',
        senderpassword: 'dlifgkgwatyidrcz',
      },
      statisticsOptions: {
        numberOfPointsHold: 5000,
      },
      decisionMakerOptions: {
        maximumBuyAmount: 0.0,
        maximumFIATBuyAmount: 500.0,
        minimumProfitPerSold: 0.00999999977648258,
        minimumPriceDropToBuy: 0.00999999977648258,
        growthDecreaseLimit: -100.0,
        growthIncreaseLimit: 100.0,
      },
      collectorOptions: {
        priceVariationDetection: 0.00999999977648258,
        datasource: null,
        newPriceTimeRate: 15,
      },
    },
    createdAt: new Date(2020, 1, 2),
  },
];

describe('ApplicationsList', () => {
  it('should render', () => {
    render(<ApplicationsList
      applications={applications}
      selectApplication={jest.fn()}
      deleteApplication={jest.fn()}
    />);
  });
});
