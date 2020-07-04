import { useEffect, useState } from 'react';

import fillDatesGaps from '../../formatters/fillDatesGaps';
import formatAssetPrices from '../../formatters/formatAssetPrices';
import formatDateTime from '../../formatters/formatDateTime';
import {
  Asset, Benchmark, BenchmarkOutput, DatesInterval,
} from '../../types';


type chartKeys = 'balances' | 'buys' | 'sells'

function filterBenchmarkResultByTime(data: BenchmarkOutput, start: number | undefined, end: number | undefined) {
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

export default (benchmark: Benchmark) => {
  const [chartDatesInterval, setChartDatesInterval] = useState<DatesInterval>();

  const [prices, setPrices] = useState<[number, number][]>();
  const [assets, setAssets] = useState<Asset[]>();
  const [balances, setBalances] = useState<[number, number][]>();
  const [buys, setBuys] = useState<[number, number][]>();
  const [sells, setSells] = useState<[number, number][]>();
  const [growth, setGrowth] = useState<[number, number][]>();
  const [growthOfGrowth, setGrowthOfGrowth] = useState<[number, number][]>();

  useEffect(() => {
    if (benchmark && benchmark.output.buys && benchmark.output.sells) {
      const sellsDates = benchmark.output.sells.map((sell) => sell[0]);
      const buysDates = benchmark.output.buys.map((buy) => buy[0]);
      const firstBuy = Math.min(...buysDates);
      const firstSell = Math.min(...sellsDates);
      const lastBuy = Math.max(...buysDates);
      const lastSell = Math.max(...sellsDates);


      const startDate = firstBuy > firstSell ? firstSell : firstBuy;
      const endDate = lastBuy > lastSell ? lastBuy : lastSell;

      if (startDate && endDate) {
        setChartDatesInterval({ startDate: new Date(startDate), endDate: new Date(endDate) });
      }
    }

    if (benchmark) {
      setBalances(
        fillDatesGaps(
          benchmark.output.balances,
        ),
      );
      setBuys(benchmark.output.buys);
      setSells(benchmark.output.sells);
      setAssets(benchmark.output.assets);
    }
  }, [benchmark]);

  useEffect(() => {
    if (chartDatesInterval) {
      const startDateFormatted = formatDateTime(chartDatesInterval.startDate.getTime());
      const endDateFormatted = formatDateTime(chartDatesInterval.endDate.getTime());

      if (benchmark) {
        fetch(`/api/assets/${benchmark?.input.asset}/prices?startDate=${startDateFormatted}&endDate=${endDateFormatted}`)
          .then((res) => res.json())
          .then(formatAssetPrices)
          .then((data) => data.sort((a, b) => b[0] - a[0]))
          .then((prices) => {
            setPrices(prices);
            const growth = derivate(prices.reverse());
            setGrowth(growth);
            setGrowthOfGrowth(derivate(growth));
          });

        const data = filterBenchmarkResultByTime(
          benchmark.output,
          chartDatesInterval.startDate.getTime(),
          chartDatesInterval.endDate.getTime(),
        );

        setBalances(
          fillDatesGaps(
            data.balances,
          ),
        );
        setBuys(data.buys);
        setSells(data.sells);
      }
    }
  }, [chartDatesInterval, benchmark]);


  return {
    prices,
    buys,
    sells,
    growth,
    growthOfGrowth,
    balances,
    assets,
    chartDatesInterval,
    setChartDatesInterval,
  };
};
