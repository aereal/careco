'use client';

import { buildMonthlyReportURL } from '@/build-monthly-report-url';
import { format } from 'date-fns/format';
import { ChangeEventHandler, FC } from 'react';

export const SelectMonth: FC = () => {
  const handleChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    e.preventDefault();
    if (!e.target.valueAsDate) {
      return;
    }
    window.location.pathname = buildMonthlyReportURL(e.target.valueAsDate);
  };
  const currentMonth = format(new Date(), 'yyyy-MM');
  return <input type='month' max={currentMonth} onChange={handleChange} />;
};
