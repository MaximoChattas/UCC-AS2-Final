import React, { useState } from "react";
import { DateRange } from "react-date-range";
import "react-date-range/dist/styles.css";
import "react-date-range/dist/theme/default.css";

const Calendar = ({ onSelectDates }) => {
  const [range, setRange] = useState([
    {
      startDate: new Date(),
      endDate: new Date(),
      key: "selection",
    },
  ]);

  const handleSelectRange = (ranges) => {
    setRange([ranges.selection]);
    onSelectDates(ranges.selection);
  };

  return (
    <div className="calendarWrap">
      <DateRange
        ranges={range}
        onChange={handleSelectRange}
        editableDateInputs={true}
        moveRangeOnFirstSelection={false}
        months={2}
        minDate={new Date()}
        direction="horizontal"
        className="calendarElement"
      />
    </div>
  );
};

export default Calendar;