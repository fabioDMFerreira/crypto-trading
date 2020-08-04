import React from 'react';

import Chart from '../components/Chart';
import ChartFilters from '../components/ChartFilters';
import useAssetPrices from '../hooks/useAssetPrices';
import useDatesInterval from '../hooks/useDatesInterval';
import useApplicationChart from './useApplicationChart';

interface Props {
  asset: string
}

export default ({ asset }: Props) => {
  const {
    datesInterval,
    setDatesInterval,
  } = useDatesInterval();

  const {
    assetPrices: prices,
  } = useAssetPrices(asset, datesInterval?.startDate, datesInterval?.endDate);

  const {
    sells,
    buys,
    balances,
    applicationState,
  } = useApplicationChart();

  return (
    <>
      <div className="mt-4">
        <ChartFilters
          minimumDate={new Date(2020, 6, 1)}
          maximumDate={new Date()}
          startDate={datesInterval?.startDate}
          endDate={datesInterval?.endDate}
          setDatesInterval={setDatesInterval}
        />
      </div>
      {
        prices && balances && sells && buys
        && (
          <div className="mt-4">
            <Chart
              prices={prices}
              balances={balances}
              buys={buys}
              sells={sells}
              setDatesInterval={setDatesInterval}
              applicationState={applicationState}
            />
          </div>
        )
      }
    </>
  );
};
