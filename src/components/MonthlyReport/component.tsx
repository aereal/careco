import { FC } from 'react';

export interface MonthlyReportProps {
  readonly year: string;
  readonly month: string;
}

export const MonthlyReport: FC<MonthlyReportProps> = ({ year, month }) => {
  return (
    <>
      <h1>
        {year}年{month}月
      </h1>
    </>
  );
};
