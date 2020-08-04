import React, { useState } from 'react';

import AssetsTable from '../../components/AssetsTable';
import Chart from '../../components/Chart';
import JsonDisplayer from '../../components/JsonDisplayer';
import { Benchmark } from '../../types';
import BenchmarkFilters from './BenchmarkFilters';
import PricesAnalysisTable from './PricesAnalysisTable';
import PricesStatisticsAnalysisTable from './PricesStatisticsAnalysisTable';
import useBenchmark from './useBenchmark';


interface Props {
  benchmark: Benchmark
}

export default ({ benchmark }: Props) => {
  const [tableView, setTableView] = useState<string>('assets');

  const {
    assets,
    prices,
    balances,
    buys,
    sells,
    growth,
    growthOfGrowth,
    chartDatesInterval,
    setChartDatesInterval,
    applicationState,
  } = useBenchmark(benchmark);

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
        <BenchmarkFilters
          minimumDate={
            benchmark.output.balances && benchmark.output.balances.length
              ? new Date(benchmark.output.balances[0][0]) : undefined
          }
          maximumDate={
            benchmark.output.balances && benchmark.output.balances.length
              ? new Date(benchmark.output.balances[benchmark.output.balances.length - 1][0]) : undefined
          }
          startDate={chartDatesInterval?.startDate}
          endDate={chartDatesInterval?.endDate}
          setDatesInterval={setChartDatesInterval}
        />
      </div>
      <div className="mb-5">
        {
          prices && balances && sells && buys && assets && growth && growthOfGrowth
          && (
            <Chart
              prices={prices}
              balances={balances}
              buys={buys}
              sells={sells}
              growth={growth}
              growthOfGrowth={growthOfGrowth}
              setDatesInterval={setChartDatesInterval}
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
        {
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
        }
      </div>
    </>
  );
};
