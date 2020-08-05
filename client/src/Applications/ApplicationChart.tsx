import React from 'react';

import Chart from '../components/Chart';
import ChartFilters from '../components/ChartFilters';
import useAssetPrices from '../hooks/useAssetPrices';
import useDatesInterval from '../hooks/useDatesInterval';
import useAccountBuysAndSells from './useAccountBuysAndSells';
import useApplicationState from './useApplicationState';


interface Props {
  asset: string,
  appID: string,
  accountID: string,
}

export default ({ asset, appID, accountID }: Props) => {
  const {
    datesInterval,
    setDatesInterval,
  } = useDatesInterval();

  const {
    assetPrices: prices,
  } = useAssetPrices(asset, datesInterval?.startDate, datesInterval?.endDate);

  const {
    applicationState,
  } = useApplicationState(appID, datesInterval?.startDate, datesInterval?.endDate);

  const {
    buys,
    sells,
  } = useAccountBuysAndSells(accountID);

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
        prices && sells && buys
        && (
          <div className="mt-4">
            <Chart
              prices={prices}
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
