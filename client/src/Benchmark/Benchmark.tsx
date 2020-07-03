import axios from 'axios';
import React, { useEffect, useState } from 'react';

import fillDatesGaps from '../formatters/fillDatesGaps';
import formatAssetPrices from '../formatters/formatAssetPrices';
import formatDateYYYYMMDD from '../formatters/formatDateYYYYMMDD';
import { Asset, BenchmarkResult } from '../types';
import AssetsTable from './AssetsTable';
import BenchmarkChart from './BenchmarkChart';
import BenchmarkFilters from './BenchmarkFilters';
import BenchmarkForm from './BenchmarkForm';
import BenchmarkList from './BenchmarksList';
import PricesAnalysisTable from './PricesAnalysisTable';
import PricesStatisticsAnalysisTable from './PricesStatisticsAnalysisTable';

type chartKeys = 'prices' | 'balances' | 'buys' | 'sells'

function filterBenchmarkResultByTime(data: BenchmarkResult, start: number | undefined, end: number | undefined) {
  const keys: chartKeys[] = ['balances', 'buys', 'sells'];
  const newData: any = {};
  keys.forEach((key) => {
    if (!data[key]) {
      return;
    }

    newData[key] = [];
    if (data[key].length && start) {
      for (let i = 0; i < data[key].length; i++) {
        if (data[key][i][0] > start) {
          newData[key] = [...data[key].slice(i)];
          break;
        }
      }
    }

    if (data[key].length && end) {
      for (let i = data[key].length - 1; i >= 0; i--) {
        if (data[key][i][0] < end) {
          newData[key] = [...newData[key].slice(0, i + 1)];
          break;
        }
      }
    }
  });

  return newData;
}

const derivate = (ns: [number, number][]): [number, number][] => {
  if (!ns.length) {
    return [];
  }
  if (ns.length === 1) {
    return ns;
  }

  return ns.slice(1).reduce((final: [number, number][], currentN: [number, number], index) => ([
    ...final,
    [currentN[0], currentN[1] - ns[index][1]],
  ]), []);
};


export default () => {
  const [benchmarkResult, setBenchmarkResult] = useState<BenchmarkResult>();
  const [prices, setPrices] = useState<[number, number][]>();
  const [assets, setAssets] = useState<Asset[]>();
  const [balances, setBalances] = useState<[number, number][]>();
  const [buys, setBuys] = useState<[number, number][]>();
  const [sells, setSells] = useState<[number, number][]>();
  const [growth, setGrowth] = useState<[number, number][]>();
  const [growthOfGrowth, setGrowthOfGrowth] = useState<[number, number][]>();

  const [tableView, setTableView] = useState<string>('assets');

  const [dataSourceOptions, setDataSourceOptions] = useState<string[]>([]);
  const [benchmarks, setBenchmarks] = useState<any>([]);

  const [startDate, setStartDate] = useState<Date>();
  const [endDate, setEndDate] = useState<Date>();

  useEffect(() => {
    fetch('/api/benchmark/data-sources')
      .then((res) => res.json())
      .then(setDataSourceOptions);

    fetch('/api/benchmark')
      .then((res) => res.json())
      .then((benchmarks) => setBenchmarks(benchmarks || []));
  }, []);

  useEffect(() => {
    if (benchmarkResult && benchmarkResult.buys && benchmarkResult.sells) {
      const sellsDates = benchmarkResult.sells.map((sell) => sell[0]);
      const buysDates = benchmarkResult.buys.map((buy) => buy[0]);
      const firstBuy = Math.min(...buysDates);
      const firstSell = Math.min(...sellsDates);
      const lastBuy = Math.max(...buysDates);
      const lastSell = Math.max(...sellsDates);


      const startDate = firstBuy > firstSell ? firstSell : firstBuy;
      const endDate = lastBuy > lastSell ? lastBuy : lastSell;

      if (startDate) {
        const date = formatDateYYYYMMDD(startDate - 10 * 24 * 60 * 60 * 1000);
        setStartDate(new Date(date));
      }
      if (endDate) {
        const date = formatDateYYYYMMDD(endDate + 10 * 24 * 60 * 60 * 1000);
        setEndDate(new Date(date));
      }
    }

    if (benchmarkResult) {
      setBalances(
        fillDatesGaps(
          benchmarkResult.balances,
        ),
      );
      setBuys(benchmarkResult.buys);
      setSells(benchmarkResult.sells);
      setAssets(benchmarkResult.assets);
    }
  }, [benchmarkResult]);

  useEffect(() => {
    if (startDate && endDate) {
      const startDateFormatted = formatDateYYYYMMDD(startDate.getTime());
      const endDateFormatted = formatDateYYYYMMDD(endDate.getTime());

      fetch(`/api/assets/BTC/prices?startDate=${startDateFormatted}&endDate=${endDateFormatted}`)
        .then((res) => res.json())
        .then(formatAssetPrices)
        .then((data) => data.sort((a, b) => b[0] - a[0]))
        .then((prices) => {
          setPrices(prices);
          const growth = derivate(prices.reverse());
          setGrowth(growth);
          setGrowthOfGrowth(derivate(growth));
        });

      if (benchmarkResult) {
        const data = filterBenchmarkResultByTime(benchmarkResult, new Date(startDate).getTime(), new Date(endDate).getTime());

        setBalances(
          fillDatesGaps(
            data.balances,
          ),
        );
        setBuys(data.buys);
        setSells(data.sells);
      }
    }
  }, [startDate, endDate, benchmarkResult]);

  // useEffect(() => {
  //   let newData: BenchmarkResult | undefined;
  //   if (benchmarkResult) {
  //     newData = JSON.parse(JSON.stringify(benchmarkResult));
  //     let start;
  //     let end;

  //     if (startDate) {
  //       start = new Date(startDate).getTime();
  //     }

  //     if (endDate) {
  //       end = new Date(endDate).getTime();
  //     }

  //     if (newData) {
  //       newData = filterBenchmarkResultByTime(newData, start, end);
  //     }
  //   }
  //   setData(newData);
  // }, [startDate, endDate, benchmarkResult]);

  const executeBenchmark = (input: any) => {
    axios
      .post('/api/benchmark', input)
      .then((res) => setBenchmarks(benchmarks.concat(res.data)));
  };

  return (
    <div>
      <div className="mt-3 mb-5">
        <BenchmarkForm onSubmit={executeBenchmark} dataSourceOptions={dataSourceOptions} />
      </div>
      <div className="mt-3 mb-5">
        <BenchmarkList
          benchmarks={benchmarks}
          selectBenchmark={(id: string) => {
            const benchmark = benchmarks.find((b: any) => b._id === id);

            if (benchmark) {
              setBenchmarkResult(JSON.parse(JSON.stringify(benchmark.output)));
            }
          }}
          deleteBenchmark={(id: string) => {
            axios.delete(`/api/benchmark/${id}`)
              .then(() => {
                setBenchmarks(benchmarks.filter((b: any) => b._id !== id));
              });
          }}
        />
      </div>
      {
        benchmarkResult
        && (
          <>
            <div className="mb-5">
              <BenchmarkFilters
                minimumDate={
                  benchmarkResult.balances && benchmarkResult.balances.length
                    ? new Date(benchmarkResult.balances[0][0]) : undefined
                }
                maximumDate={
                  benchmarkResult.balances && benchmarkResult.balances.length
                    ? new Date(benchmarkResult.balances[benchmarkResult.balances.length - 1][0]) : undefined
                }
                startDate={startDate}
                endDate={endDate}
                setStartDate={setStartDate}
                setEndDate={setEndDate}
              />
            </div>
            <div className="mb-5">
              {
                prices && balances && sells && buys && assets && growth && growthOfGrowth
                && (
                  <BenchmarkChart
                    prices={prices}
                    balances={balances}
                    buys={buys}
                    sells={sells}
                    growth={growth}
                    growthOfGrowth={growthOfGrowth}
                    setStartDate={setStartDate}
                    setEndDate={setEndDate}
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
        )
      }
    </div>
  );
};
