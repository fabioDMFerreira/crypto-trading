import { fireEvent, render } from '@testing-library/react';
import React from 'react';
import { act } from 'react-dom/test-utils';

import BenchmarkForm from './BenchmarkForm';

const fnMock = () => { };

const dataSourceOptionsMock = {
  btc: {
    'Last Year Minute': 'btc/last-year-minute.csv',
  },
  eth: {
    'Last Year Minute': 'eth/last-year-minute.csv',
  },
};

describe('BenchmarkForm', () => {
  it('should render', () => {
    render(<BenchmarkForm dataSourceOptions={dataSourceOptionsMock} onSubmit={fnMock} />);
  });

  it('should call onSubmit on clicking submit button', async () => {
    const onSubmit = jest.fn();
    const { getByTestId } = render(<BenchmarkForm dataSourceOptions={dataSourceOptionsMock} onSubmit={onSubmit} />);

    const submitButton = getByTestId('submit-button');

    expect(submitButton).toBeTruthy();

    await act(async () => {
      fireEvent.click(submitButton);
    });

    expect(onSubmit).toHaveBeenCalledTimes(1);

    expect(onSubmit.mock.calls[0][0]).toEqual({
      decisionMakerOptions: {
        maximumBuyAmount: 0.1,
        minimumProfitPerSold: 0.03,
        minimumPriceDropToBuy: 0.01,
      },
      statisticsOptions: {
        numberOfPointsHold: 2000,
      },
      collectorOptions: {
        priceVariationDetection: 0.01,
      },
      accountInitialAmount: 2000,
      dataSourceFilePath: 'btc/last-year-minute.csv',
      asset: 'btc',
    });
  });
});
