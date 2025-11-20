import {
  fragmentTotalDistance,
  TotalDrivingDistance,
} from '@/components/TotalDrivingDistance';
import { FragmentType, getFragmentData } from '@/graphql';
import { numberOf } from '@/month';
import { FC } from 'react';
import { fragmentChartDataSeries, RecordChart } from '../RecordChart';
import { fragmentMonthlySummary } from './fragment.monthly-report';

type Data = FragmentType<typeof fragmentMonthlySummary> &
  FragmentType<typeof fragmentTotalDistance>;

interface MonthlyReportProps {
  readonly summary: Data;
  readonly records: FragmentType<typeof fragmentChartDataSeries>;
}

export const MonthlyReport: FC<MonthlyReportProps> = ({ summary, records }) => {
  const { year, month } = getFragmentData(fragmentMonthlySummary, summary);
  return (
    <div className='h-dvh'>
      <h1 className='font-bold text-2xl mt-4 -mb-4'>
        {year}年{numberOf(month)}月の走行記録
      </h1>
      <TotalDrivingDistance {...summary} />
      <div className='my-8 h-[70vh]'>
        <RecordChart records={records} omitMonth />
      </div>
    </div>
  );
};
