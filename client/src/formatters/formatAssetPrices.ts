export interface AssetPriceAggregatedByDate {
  _id: {
    day: number;
    month: number;
    year: number;
    hour?: number;
    minute?: number;
  };
  price: number;
}

export default (data: AssetPriceAggregatedByDate[]): [number, number][] => (
  data
    ? data
      .map(
        (assetPrice) => [
          Date.UTC(assetPrice._id.year, assetPrice._id.month - 1, assetPrice._id.day, assetPrice._id.hour || 0, assetPrice._id.minute || 0),
          assetPrice.price,
        ],
      )
    : []
);
