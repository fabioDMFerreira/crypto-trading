import { useEffect, useState } from 'react';

import formatApplicationExecutionState, { ApplicationState } from '../formatters/formatApplicationExecutionState';
import formatDateTime from '../formatters/formatDateTime';

export default (appID: string, startDate: Date, endDate: Date) => {
  const [applicationState, setApplicationState] = useState<ApplicationState>();

  useEffect(() => {
    if (appID && startDate && endDate) {
      const startDateFormatted = formatDateTime(startDate.getTime());
      const endDateFormatted = formatDateTime(endDate.getTime());

      fetch(`/api/applications/${appID}/state?startDate=${startDateFormatted}&endDate=${endDateFormatted}`)
        .then((res) => res.json())
        .then(formatApplicationExecutionState)
        .then((data) => ({
          average: data.average.sort((a, b) => a[0] - b[0]),
          standardDeviation: data.standardDeviation.sort((a, b) => a[0] - b[0]),
          currentChange: data.currentChange.sort((a, b) => a[0] - b[0]),
          lowerBollingerBand: data.lowerBollingerBand.sort((a, b) => a[0] - b[0]),
          higherBollingerBand: data.higherBollingerBand.sort((a, b) => a[0] - b[0]),
          accountAmount: data.accountAmount.sort((a, b) => a[0] - b[0]),
        }))
        .then(setApplicationState);
    } else {
      setApplicationState(undefined);
    }
  }, [appID, startDate, endDate]);

  return {
    applicationState,
  };
};
