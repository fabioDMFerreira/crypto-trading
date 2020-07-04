import * as Highcharts from 'highcharts';
import HighchartsReact from 'highcharts-react-official';
import React from 'react';
import useDebouncedCallback from 'use-debounce/lib/useDebouncedCallback';

// const getDates = (el: [number, number]) => new Date(el[0]);

interface BenchmarkChartProps {
  prices: [number, number][]
  balances: [number, number][]
  buys: [number, number][]
  sells: [number, number][]
  growth: [number, number][]
  growthOfGrowth: [number, number][]
  setStartDate: (startDate: Date) => void
  setEndDate: (endDate: Date) => void
}

export default ({
  prices,
  balances,
  buys,
  sells,
  growth,
  growthOfGrowth,
  setStartDate, setEndDate,
}: BenchmarkChartProps) => {
  const [setRange] = useDebouncedCallback((min: number, max: number) => {
    setStartDate(new Date(min));
    setEndDate(new Date(max));
  }, 100);

  // console.log({
  //   prices: prices.map(getDates),
  //   balances: balances.map(getDates),
  //   buys: buys.map(getDates),
  //   sells: sells.map(getDates),
  //   growth: growth.map(getDates),
  //   growthOfGrowth: growthOfGrowth.map(getDates),
  // });

  const options: Highcharts.Options = {
    title: {
      text: 'Benchmark',
    },
    chart: {
      zoomType: 'x',
    },
    xAxis: {
      type: 'datetime',
      dateTimeLabelFormats: {
        day: '%e-%b-%y',
        month: '%b-%y',
      },
      labels: {
        formatter() {
          return Highcharts.dateFormat('%d-%b-%y', (this.value));
        },
      },
      events: {
        afterSetExtremes: (e) => {
          if (e.trigger === 'zoom') {
            setRange(e.min, e.max);
          }
        },
      },
    },
    yAxis: [
      {},
      { opposite: true },
      { opposite: true },
    ],
    tooltip: {
      shared: true,
    },
    series: [{
      name: 'Prices',
      type: 'line',
      data: prices,
      color: 'rgba(83, 83, 223, .5)',
    }, {
      name: 'Growth',
      type: 'line',
      data: growth,
      yAxis: 2,
      visible: false,
      color: '#FFA000',
    }, {
      name: 'Growth of growth',
      type: 'line',
      data: growthOfGrowth,
      yAxis: 2,
      visible: false,
      color: '#FFF176',
    }, {
      name: 'Balance',
      type: 'line',
      yAxis: 1,
      data: balances,
      color: 'rgba(83, 223, 223, .5)',
    }, {
      name: 'Buys',
      type: 'scatter',
      data: buys,
      color: 'rgba(83, 223, 83, .5)',
      marker: {
        radius: 5,
        symbol: 'circle',
      },
    }, {
      name: 'Sells',
      type: 'scatter',
      data: sells,
      color: 'rgba(223, 83, 83, .5)',
      marker: {
        radius: 5,
        symbol: 'circle',
      },
    }],
  };

  return (
    <HighchartsReact
      highcharts={Highcharts}
      options={options}
    />
  );
};
