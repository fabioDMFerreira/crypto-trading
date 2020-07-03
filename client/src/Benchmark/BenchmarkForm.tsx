import React from 'react';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import { useForm } from 'react-hook-form';

interface DecisionMakerOptions {
  maximumBuyAmount: number,
  minimumProfitPerSold: number,
  minimumPriceDropToBuy: number
}

interface StatisticsOptions {
  numberOfPointsHold: number
}

interface CollectorOptions {
  priceVariationDetection: number
}

interface BenchmarkInput {
  decisionMakerOptions: DecisionMakerOptions,
  statisticsOptions: StatisticsOptions
  collectorOptions: CollectorOptions
  accountInitialAmount: number,
  dataSourceFilePath: string,
}

interface Props {
  onSubmit: (data: any) => void,
  dataSourceOptions: string[]
}

const benchmarkDefaults: BenchmarkInput = {
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
};

const serializeBenchmarkInput = (fn: any) => (data: any) => {
  const dataSerialized = {
    decisionMakerOptions: {
      maximumBuyAmount: +data.decisionMakerOptions.maximumBuyAmount,
      minimumProfitPerSold: +data.decisionMakerOptions.minimumProfitPerSold,
      minimumPriceDropToBuy: +data.decisionMakerOptions.minimumPriceDropToBuy,
    },

    statisticsOptions: {
      numberOfPointsHold: +data.statisticsOptions.numberOfPointsHold,
    },

    collectorOptions: {
      priceVariationDetection: +data.collectorOptions.priceVariationDetection,
    },

    dataSourceFilePath: data.dataSourceFilePath,

    accountInitialAmount: +data.accountInitialAmount,
  };

  fn(dataSerialized);
};

export default ({ onSubmit, dataSourceOptions }: Props) => {
  const { register, handleSubmit } = useForm();

  return (
    <Form onSubmit={handleSubmit(serializeBenchmarkInput(onSubmit))}>
      <Form.Row>
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

        <Form.Group as={Col} controlId="formGridDataSourceFilePath">
          <Form.Label>Data Source File Path</Form.Label>
          <Form.Control
            defaultValue={benchmarkDefaults.dataSourceFilePath}
            name="dataSourceFilePath"
            as="select"
            placeholder="Enter file path"
            ref={register}
          >
            {
              dataSourceOptions.map(
                (dataSource) => (
                  <option>{dataSource}</option>
                ),
              )
            }
          </Form.Control>
        </Form.Group>
      </Form.Row>

      <Button variant="primary" type="submit" data-testid="submit-button">
        Benchmark
      </Button>
    </Form>
  );
};
