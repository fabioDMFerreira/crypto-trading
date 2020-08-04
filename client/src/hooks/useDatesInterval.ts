import { useState } from 'react';

import { DatesInterval } from '../types';

function getLastWeek() {
  const today = new Date();
  const lastWeek = new Date(today.getFullYear(), today.getMonth(), today.getDate() - 7);
  return lastWeek;
}

export default () => {
  const [datesInterval, setDatesInterval] = useState<DatesInterval>({
    startDate: getLastWeek(),
    endDate: new Date(),
  });

  return {
    datesInterval,
    setDatesInterval,
  };
};
