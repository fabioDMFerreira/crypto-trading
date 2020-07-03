const oneDayMiliseconds = 60 * 60 * 24 * 1000;

export default (values: [number, number][]) => {
  if (!values.length) {
    return [];
  } if (values.length === 1) {
    return values;
  }

  let [currentDate, currentValue] = values[0];
  const newValues = [values[0]];

  for (let i = 1; i < values.length; i++) {
    const cur = values[i];

    while (cur[0] - currentDate > oneDayMiliseconds) {
      currentDate += oneDayMiliseconds;
      newValues.push([currentDate, currentValue]);
    }

    [currentDate, currentValue] = cur;

    newValues.push(cur);
  }

  return newValues;
};
