export default function (date: number | string) {
  const d = new Date(date);
  let month = `${d.getUTCMonth() + 1}`;
  let day = `${d.getUTCDate()}`;
  const year = d.getUTCFullYear();
  const hours = d.getUTCHours();
  const minutes = d.getUTCMinutes();
  const seconds = d.getUTCSeconds();

  if (month.length < 2) { month = `0${month}`; }
  if (day.length < 2) { day = `0${day}`; }

  return `${[year, month, day].join('-')}T${[hours, minutes, seconds].join(':')}`;
}
