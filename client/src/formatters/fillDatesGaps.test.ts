import fillDatesGaps from './fillDatesGaps';

describe('fillDatesGaps', () => {
  it('should fill gap values with the previous value', () => {
    const actual = fillDatesGaps([
      [
        new Date('2020-03-01').getTime(),
        1000,
      ], [
        new Date('2020-03-05').getTime(),
        2000,
      ], [
        new Date('2020-03-07').getTime(),
        3000,
      ],
    ]);

    const expected = [
      [
        new Date('2020-03-01').getTime(),
        1000,
      ], [
        new Date('2020-03-02').getTime(),
        1000,
      ],
      [
        new Date('2020-03-03').getTime(),
        1000,
      ], [
        new Date('2020-03-04').getTime(),
        1000,
      ], [
        new Date('2020-03-05').getTime(),
        2000,
      ], [
        new Date('2020-03-06').getTime(),
        2000,
      ], [
        new Date('2020-03-07').getTime(),
        3000,
      ],
    ];

    expect(actual).toEqual(expected);
  });
});
