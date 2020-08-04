import moment from 'moment';
import TimePicker from 'rc-time-picker';
import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';

import formatDateYYYYMMDD from '../../formatters/formatDateYYYYMMDD';
import { DatesInterval } from '../../types';

interface Props {
  minimumDate: Date | undefined,
  maximumDate: Date | undefined,
  startDate: Date | undefined,
  endDate: Date | undefined,
  setDatesInterval: (interval: DatesInterval) => void,
}

export default ({
  minimumDate,
  maximumDate,
  startDate,
  endDate,
  setDatesInterval,
}: Props) => {
  const modifiers = {
    from: startDate,
    to: endDate,
  };
  // const to = useRef<DayPickerInput>();

  return (
    <div className="mb-4">
      <div className="InputFromTo">
        {
          minimumDate
          && (
            <button
              type="button"
              onClick={() => endDate && setDatesInterval({
                startDate: minimumDate,
                endDate,
              })}
            >
              {formatDateYYYYMMDD(minimumDate.getTime())}
            </button>
          )
        }
        {' '}
        <DayPickerInput
          value={startDate}
          placeholder="From"
          // format="LL"
          // formatDate={formatDate}
          // parseDate={parseDate}
          dayPickerProps={{
            selectedDays: startDate && endDate ? [startDate, { from: startDate, to: endDate }] : undefined,
            disabledDays: endDate ? { after: endDate } : undefined,
            toMonth: endDate || maximumDate,
            modifiers,
            numberOfMonths: 1,
            // onDayClick() {
            //   if (to && to.current) {
            //     to.current.getInput().focus();
            //   }
            // },
          }}
          onDayChange={(startDate) => endDate && setDatesInterval({
            startDate,
            endDate,
          })}
        />
        <TimePicker
          showSecond={false}
          value={moment(startDate?.getTime())}
          onChange={(data: moment.Moment) => {
            if (endDate && startDate) {
              setDatesInterval({
                startDate: data.toDate(),
                endDate,
              });
            }
          }}
          use12Hours
          inputReadOnly
        />
        {' - '}
        <span className="InputFromTo-to">
          <DayPickerInput
            // // @ts-ignore:
            // ref={to}
            value={endDate}
            placeholder="To"
            // format="LL"
            // formatDate={formatDate}
            // parseDate={parseDate}
            dayPickerProps={{
              selectedDays: startDate && endDate ? [startDate, { from: startDate, to: endDate }] : undefined,
              disabledDays: startDate ? { before: startDate } : undefined,
              modifiers,
              month: startDate,
              fromMonth: startDate || minimumDate,
              numberOfMonths: 1,
            }}
            onDayChange={(endDate) => startDate && setDatesInterval({
              startDate,
              endDate,
            })}
          />
          <TimePicker
            showSecond={false}
            value={moment(endDate?.getTime())}
            onChange={(data: moment.Moment) => {
              if (endDate && startDate) {
                setDatesInterval({
                  startDate,
                  endDate: data.toDate(),
                });
              }
            }}
            use12Hours
            inputReadOnly
          />
        </span>
        {' '}
        {
          maximumDate
          && (
            <button
              type="button"
              onClick={() => startDate && setDatesInterval({
                startDate,
                endDate: maximumDate,
              })}
            >
              {formatDateYYYYMMDD(maximumDate.getTime())}
            </button>
          )
        }
      </div>
    </div>
  );
};
