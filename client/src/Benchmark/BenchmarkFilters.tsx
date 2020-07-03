import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';

import formatDateYYYYMMDD from '../formatters/formatDateYYYYMMDD';

interface Props {
  minimumDate: Date | undefined,
  maximumDate: Date | undefined,
  startDate: Date | undefined,
  endDate: Date | undefined,
  setStartDate: (startDate: Date) => void,
  setEndDate: (endDate: Date) => void,
}

export default ({
  minimumDate,
  maximumDate,
  startDate,
  endDate,
  setStartDate,
  setEndDate,
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
              onClick={() => setStartDate(minimumDate)}
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
          onDayChange={setStartDate}
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
            onDayChange={setEndDate}
          />
        </span>
        {' '}
        {
          maximumDate
          && (
            <button
              type="button"
              onClick={() => setEndDate(maximumDate)}
            >
              {formatDateYYYYMMDD(maximumDate.getTime())}
            </button>
          )
        }
      </div>
    </div>
  );
};
