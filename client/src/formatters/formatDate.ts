export default function (date: number | string) {
  const d = new Date(+date);
  let month = `${d.getMonth() + 1}`;
  let day = `${d.getDate()}`;
  const year = d.getFullYear();
  const hour = d.getHours();
  const minute = d.getMinutes();

  if (month.length < 2) { month = `0${month}`; }
  if (day.length < 2) { day = `0${day}`; }

  return `${[year, month, day].join('-')} ${hour}:${minute}`;
}
