import { useEffect, useState } from 'react';

import formatAssetPrices from '../formatters/formatAssetPrices';
import formatDateTime from '../formatters/formatDateTime';

export default (asset: string, startDate: Date | undefined, endDate: Date | undefined) => {
  const [assetPrices, setAssetPrices] = useState<[number, number][]>();

  useEffect(() => {
    if (!startDate || !endDate || !asset) {
      setAssetPrices([]);
      return;
    }

    const startDateFormatted = formatDateTime(startDate.getTime());
    const endDateFormatted = formatDateTime(endDate.getTime());

    fetch(`/api/assets/${asset}/prices?startDate=${startDateFormatted}&endDate=${endDateFormatted}`)
      .then((res) => res.json())
      .then(formatAssetPrices)
      .then((data) => data.sort((a, b) => a[0] - b[0]))
      .then(setAssetPrices);
  }, [asset, startDate, endDate]);

  return {
    assetPrices,
  };
};
