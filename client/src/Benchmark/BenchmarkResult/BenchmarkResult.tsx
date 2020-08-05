import React, { useState } from 'react';

import AssetsTable from '../../components/AssetsTable';
import Chart from '../../components/Chart';
import ChartFilters from '../../components/ChartFilters/ChartFilters';
import JsonDisplayer from '../../components/JsonDisplayer';
import useAssetPrices from '../../hooks/useAssetPrices';
import useDatesInterval from '../../hooks/useDatesInterval';
import { Benchmark } from '../../types';
// import PricesAnalysisTable from './PricesAnalysisTable';
// import PricesStatisticsAnalysisTable from './PricesStatisticsAnalysisTable';
import useBenchmark from './useBenchmark';

interface Props {
  benchmark: Benchmark
}

export default ({ benchmark }: Props) => {
  const [tableView, setTableView] = useState<string>('assets');

  const {
    datesInterval,
    setDatesInterval,
    minDate,
    setMinDate,
    maxDate,
    setMaxDate,
  } = useDatesInterval();

  const {
    assets,
    buys,
    sells,
    applicationState,
  } = useBenchmark(benchmark, datesInterval, setDatesInterval, setMinDate, setMaxDate);

  const {
    assetPrices: prices,
  } = useAssetPrices(benchmark.input.asset, datesInterval?.startDate, datesInterval?.endDate);

  return (
    <>
      <div className="mb-5">
        <h2>Benchmark Result</h2>
        <div className="mb-5">
          <JsonDisplayer json={{
            ...benchmark,
            output: {
              finalAmount: benchmark.output.finalAmount,
            },
          }}
          />
        </div>
        <ChartFilters
          minimumDate={
            minDate
          }
          maximumDate={
            maxDate
          }
          startDate={datesInterval?.startDate}
          endDate={datesInterval?.endDate}
          setDatesInterval={setDatesInterval}
        />
      </div>
      <div className="mb-5">
        {
          prices && sells && buys && applicationState
          && (
            <Chart
              prices={prices}
              buys={buys}
              sells={sells}
              setDatesInterval={setDatesInterval}
              applicationState={applicationState}
            />
          )
        }
      </div>
      <div className="mb-5">
        <button type="button" onClick={() => { setTableView('assets'); }}>Assets</button>
        <button type="button" onClick={() => { setTableView('price-analysis'); }}>Price analysis</button>
      </div>
      <div className="mb-5">
        {
          assets && tableView === 'assets'
          && <AssetsTable assets={assets} />
        }
        {/* {
          prices && growth && growthOfGrowth && tableView === 'price-analysis'
          && (
            <PricesStatisticsAnalysisTable
              prices={prices}
              growth={growth}
              growthOfGrowth={growthOfGrowth}
            />
          )
        }
        {
          prices && growth && growthOfGrowth && tableView === 'price-analysis'
          && (
            <PricesAnalysisTable
              prices={prices}
              growth={growth}
              growthOfGrowth={growthOfGrowth}
            />
          )
        } */}
      </div>
    </>
  );
};
