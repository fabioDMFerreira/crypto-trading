export default function (date: number | string) {
  const d = new Date(date);
  let month = `${d.getUTCMonth() + 1}`;
  let day = `${d.getUTCDate()}`;
  const year = d.getUTCFullYear();
  let hours = `${d.getUTCHours()}`;
  let minutes = `${d.getUTCMinutes()}`;
  let seconds = `${d.getUTCSeconds()}`;

  if (month.length < 2) { month = `0${month}`; }
  if (day.length < 2) { day = `0${day}`; }
  if (hours.length < 2) { hours = `0${hours}`; }
  if (minutes.length < 2) { minutes = `0${minutes}`; }
  if (seconds.length < 2) { seconds = `0${seconds}`; }


  return `${[year, month, day].join('-')}T${[hours, minutes, seconds].join(':')}`;
}
