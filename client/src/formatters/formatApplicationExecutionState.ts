export interface ApplicationStateAggregatedByDate {
  _id: {
    day: number;
    month: number;
    year: number;
    hour?: number;
    minute?: number;
  };
  average: number;
  standardDeviation: number;
  currentChange: number
  lowerBollingerBand: number;
  higherBollingerBand: number;

  accountAmount:number;
}

export interface ApplicationState {
  average: [number, number][]
  standardDeviation: [number, number][]
  lowerBollingerBand: [number, number][]
  higherBollingerBand: [number, number][]
  currentChange: [number, number][]
  accountAmount: [number, number][]
}

export default (data: ApplicationStateAggregatedByDate[]): ApplicationState => {
  const defaultResult: ApplicationState = {
    average: [], standardDeviation: [], higherBollingerBand: [], lowerBollingerBand: [], currentChange: [], accountAmount: [],
  };

  return data
    ? data
      .reduce(
        (final, applicationState) => {
          const time = Date.UTC(
            applicationState._id.year,
            applicationState._id.month - 1,
            applicationState._id.day,
            applicationState._id.hour || 0,
            applicationState._id.minute || 0,
          );

          final.average.push([time, applicationState.average]);
          final.standardDeviation.push([time, applicationState.standardDeviation]);
          final.currentChange.push([time, applicationState.currentChange]);
          final.lowerBollingerBand.push([time, applicationState.lowerBollingerBand]);
          final.higherBollingerBand.push([time, applicationState.higherBollingerBand]);
          final.accountAmount.push([time, applicationState.accountAmount]);

          return final;
        }, defaultResult,
      )
    : defaultResult;
};
