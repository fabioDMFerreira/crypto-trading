import formatAssetPrices, { AssetPriceAggregatedByDate } from './formatAssetPrices';

const fixture: AssetPriceAggregatedByDate[] = [{
  _id: { day: 28, month: 5, year: 2020 }, price: 9507.5244140625,
},
{ _id: { day: 14, month: 11, year: 2019 }, price: 8656.970703125 },
{ _id: { day: 26, month: 11, year: 2019 }, price: 7161.99658203125 },
];

describe('formatAssetPrices', () => {
  it('should convert to a highchart readable format', () => {
    const actual = formatAssetPrices(fixture);
    const expected = [[1590624000000, 9507.5244140625], [1573689600000, 8656.970703125], [1574726400000, 7161.99658203125]];

    expect(actual).toEqual(expected);
  });
});
