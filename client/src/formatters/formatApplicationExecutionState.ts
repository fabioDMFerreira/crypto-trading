export interface ApplicationStateAggregatedByDate {
  _id: {
    day: number;
    month: number;
    year: number;
    hour?: number;
    minute?: number;
  };
  average: number;
  standarddeviation: number;
  currentchange: number
  lowerbollingerband: number;
  higherbollingerband: number;
}

export interface ApplicationState {
  average: [number, number][]
  standarddeviation: [number, number][]
  lowerbollingerband: [number, number][]
  higherbollingerband: [number, number][]
  currentchange: [number, number][]
}

export default (data: ApplicationStateAggregatedByDate[]): ApplicationState => {
  const defaultResult: ApplicationState = {
    average: [], standarddeviation: [], higherbollingerband: [], lowerbollingerband: [], currentchange: [],
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
          final.standarddeviation.push([time, applicationState.standarddeviation]);
          final.currentchange.push([time, applicationState.currentchange]);
          final.lowerbollingerband.push([time, applicationState.lowerbollingerband]);
          final.higherbollingerband.push([time, applicationState.higherbollingerband]);

          return final;
        }, defaultResult,
      )
    : defaultResult;
};
