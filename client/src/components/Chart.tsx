import * as Highcharts from 'highcharts';
import HighchartsReact from 'highcharts-react-official';
import React from 'react';
import useDebouncedCallback from 'use-debounce/lib/useDebouncedCallback';

import { ApplicationState } from '../formatters/formatApplicationExecutionState';
import { DatesInterval } from '../types';

interface ChartProps {
  prices: [number, number][]
  balances: [number, number][]
  buys: [number, number][]
  sells: [number, number][]
  growth: [number, number][]
  growthOfGrowth: [number, number][]
  applicationState?: ApplicationState
  setDatesInterval: (interval: DatesInterval) => void
}

export default ({
  prices,
  balances,
  buys,
  sells,
  growth,
  growthOfGrowth,
  applicationState,
  setDatesInterval,
}: ChartProps) => {
  const [setRange] = useDebouncedCallback((min: number, max: number) => {
    setDatesInterval({
      startDate: new Date(min),
      endDate: new Date(max),
    });
  }, 100);

  const options: Highcharts.Options = {
    chart: {
      zoomType: 'x',
    },
    xAxis: {
      type: 'datetime',
      dateTimeLabelFormats: {
        day: '%e-%b-%y',
        month: '%b-%y',
        second: '%H:%M:%S',
        minute: '%H:%M',
        hour: '%H:%M',
      },
      labels: {
        formatter() {
          // eslint-disable-next-line react/no-this-in-sfc
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
      xDateFormat: '%Y-%m-%d %H:%M:%S',
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
      tooltip: {
        pointFormat: 'x: <b>{point.x:%d-%m-%y %H:%M:%S}</b><br/>y: <b>{point.y}</b><br/>',
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
      tooltip: {
        pointFormat: 'x: <b>{point.x:%y-%m-%d %H:%M:%S}</b><br/>y: <b>{point.y}</b><br/>',
      },
    }],
  };

  if (applicationState && options.series) {
    options.series = options.series.concat([{
      name: 'Average',
      type: 'line',
      data: applicationState.average,
      visible: false,
      color: '#000',
    }, {
      name: 'L Bollinger',
      type: 'line',
      data: applicationState.lowerbollingerband,
      visible: false,
      color: '#ccc',
    }, {
      name: 'H Bollinger',
      type: 'line',
      data: applicationState.higherbollingerband,
      visible: false,
      color: '#ccc',
    }, {
      name: 'Standard Deviation',
      type: 'line',
      data: applicationState.standarddeviation,
      visible: false,
      color: '#FFF176',
    }, {
      name: 'Current Change',
      type: 'line',
      data: applicationState.currentchange,
      yAxis: 2,
      visible: false,
      color: '#FFF176',
    }]);
  }

  return (
    <HighchartsReact
      highcharts={Highcharts}
      options={options}
    />
  );
};
