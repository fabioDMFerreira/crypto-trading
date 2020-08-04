import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import { useForm } from 'react-hook-form';

import { BenchmarkInput, DataSourceOptions } from '../../types';
import useDataSourceOptionsParser from './useDataSourceOptionsParser';


interface Props {
  onSubmit: (data: any) => void,
  dataSourceOptions: DataSourceOptions
}

const benchmarkDefaults: BenchmarkInput = {
  decisionMakerOptions: {
    maximumBuyAmount: 0.1,
    maximumFIATBuyAmount: 0,
    minimumProfitPerSold: 0.02,
    minimumPriceDropToBuy: 0.01,
    growthIncreaseLimit: 100,
    growthDecreaseLimit: -100,
  },
  statisticsOptions: {
    numberOfPointsHold: 20000,
  },
  collectorOptions: {
    priceVariationDetection: 0.01,
    datasource: null,
    newPriceTimeRate: 15,
  },
  accountInitialAmount: 5000,
  asset: 'btc',
  dataSourceFilePath: 'btc/last-year-minute.csv',
};

const serializeBenchmarkInput = (fn: any) => (data: any) => {
  const dataSerialized = {
    decisionMakerOptions: {
      maximumBuyAmount: data.decisionMakerOptions.maximumBuyAmount ? +data.decisionMakerOptions.maximumBuyAmount : undefined,
      maximumFIATBuyAmount: data.decisionMakerOptions.maximumFIATBuyAmount ? +data.decisionMakerOptions.maximumFIATBuyAmount : undefined,
      minimumProfitPerSold: +data.decisionMakerOptions.minimumProfitPerSold,
      minimumPriceDropToBuy: +data.decisionMakerOptions.minimumPriceDropToBuy,
      growthIncreaseLimit: +data.decisionMakerOptions.growthIncreaseLimit,
      growthDecreaseLimit: +data.decisionMakerOptions.growthDecreaseLimit,
    },

    statisticsOptions: {
      numberOfPointsHold: +data.statisticsOptions.numberOfPointsHold,
    },

    collectorOptions: {
      priceVariationDetection: +data.collectorOptions.priceVariationDetection,
      newPriceTimeRate: +data.collectorOptions.newPriceTimeRate,
    },

    dataSourceFilePath: data.dataSourceFilePath,

    accountInitialAmount: +data.accountInitialAmount,
  };

  fn(dataSerialized);
};

export default ({ onSubmit, dataSourceOptions }: Props) => {
  const { register, handleSubmit } = useForm();

  const {
    assets, dataSources, activeAsset, setActiveAsset,
    activeDataSource, setActiveDataSource,
  } = useDataSourceOptionsParser(dataSourceOptions);

  const [maximumBuyAmountCurrency, setMaximumBuyAmountCurrency] = useState('Asset');

  const submit = (data: any) => {
    if (maximumBuyAmountCurrency === 'FIAT') {
      data.decisionMakerOptions.maximumFIATBuyAmount = data.decisionMakerOptions.maximumBuyAmount;
      delete data.decisionMakerOptions.maximumBuyAmount;
    }

    onSubmit({
      ...data,
      asset: activeAsset,
    });
  };

  return (
    <Form onSubmit={handleSubmit(serializeBenchmarkInput(submit))}>
      <Form.Row>
        <Form.Group as={Col} controlId="formGridAmount">
          <Form.Label>Account Initial Amount</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.accountInitialAmount}
            name="accountInitialAmount"
            type="number"
            placeholder="Enter amount"
            ref={register}
          />
        </Form.Group>

        <Form.Group>
          <Form.Label>Maximum Buy Amount Currency</Form.Label>
          <Form.Control
            defaultValue={maximumBuyAmountCurrency}
            name="maximumBuyAmountCurrency"
            as="select"
            onChange={(e) => setMaximumBuyAmountCurrency(e.target.value)}
          >
            <option value="Asset">Asset</option>
            <option value="FIAT">FIAT</option>
          </Form.Control>
        </Form.Group>

        <Form.Group as={Col} controlId="formGridMaximumBuyAmount">
          <Form.Label>Maximum Buy Amount</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.decisionMakerOptions.maximumBuyAmount}
            name="decisionMakerOptions.maximumBuyAmount"
            type="number"
            placeholder="Enter amount"
            step="0.1"
            ref={register}
          />
        </Form.Group>

        <Form.Group as={Col} controlId="formGridMinimumProfitPerSold">
          <Form.Label>Minimum Profit Per Sold</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.decisionMakerOptions.minimumProfitPerSold}
            name="decisionMakerOptions.minimumProfitPerSold"
            type="number"
            placeholder="Enter amount"
            step="0.01"
            ref={register}
          />
        </Form.Group>

        <Form.Group as={Col} controlId="formGridMinimumPriceDropToBuy">
          <Form.Label>Minimum Price Drop To Buy</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.decisionMakerOptions.minimumPriceDropToBuy}
            name="decisionMakerOptions.minimumPriceDropToBuy"
            type="number"
            placeholder="Enter amount"
            step="0.01"
            ref={register}
          />
        </Form.Group>
      </Form.Row>

      <Form.Row>
        <Form.Group as={Col} controlId="formGridStatisticsPointsToHold">
          <Form.Label>Number of Points to Hold</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.statisticsOptions.numberOfPointsHold}
            name="statisticsOptions.numberOfPointsHold"
            type="number"
            placeholder="Enter amount"
            ref={register}
          />
        </Form.Group>
        <Form.Group as={Col} controlId="formMinutesToCollectNewPoint">
          <Form.Label>Minutes to collect new point</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.collectorOptions.newPriceTimeRate}
            name="collectorOptions.newPriceTimeRate"
            type="number"
            placeholder="Enter time in minutes"
            ref={register}
          />
        </Form.Group>

        <Form.Group as={Col} controlId="formGrowthIncreaseLimit">
          <Form.Label>Growth Increase Limit</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.decisionMakerOptions.growthIncreaseLimit}
            name="decisionMakerOptions.growthIncreaseLimit"
            type="number"
            min="0"
            placeholder="Enter increase limit"
            ref={register}
          />
        </Form.Group>

        <Form.Group as={Col} controlId="formGrowthDecreaseLimit">
          <Form.Label>Growth Decrease Limit</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.decisionMakerOptions.growthDecreaseLimit}
            name="decisionMakerOptions.growthDecreaseLimit"
            type="number"
            placeholder="Enter decrease limit"
            ref={register}
            max="0"
          />
        </Form.Group>
      </Form.Row>

      <Form.Row>
        <Form.Group as={Col} controlId="formGridPriceVariationDetection">
          <Form.Label>Price Variation Detection</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.collectorOptions.priceVariationDetection}
            name="collectorOptions.priceVariationDetection"
            type="number"
            placeholder="Enter amount"
            step="0.01"
            ref={register}
          />
        </Form.Group>


        {
          assets
          && (
            <Form.Group as={Col} controlId="formGridAsset">
              <Form.Label>Asset</Form.Label>
              <Form.Control
                value={activeAsset}
                name="asset"
                as="select"
                placeholder="Enter asset"
                ref={register}
                onChange={(e) => {
                  setActiveAsset(e.target.value);
                }}
              >
                {
                  assets.map(
                    (asset) => (
                      <option key={`${asset.value}:${asset.label}`} value={asset.value}>{asset.label}</option>
                    ),
                  )
                }
              </Form.Control>
            </Form.Group>
          )
        }

        {
          dataSources
          && (
            <Form.Group as={Col} controlId="formGridDataSourceFilePath">
              <Form.Label>Data Source File Path</Form.Label>
              <Form.Control
                value={activeDataSource}
                onChange={(e) => {
                  setActiveDataSource(e.target.value);
                }}
                name="dataSourceFilePath"
                as="select"
                placeholder="Enter file path"
                ref={register}
              >
                {
                  dataSources.map(
                    (dataSource) => (
                      <option key={`${dataSource.value}:${dataSource.label}`} value={dataSource.value}>{dataSource.label}</option>
                    ),
                  )
                }
              </Form.Control>
            </Form.Group>
          )
        }
      </Form.Row>

      <Button variant="primary" type="submit" data-testid="submit-button">
        Benchmark
      </Button>
    </Form>
  );
};
