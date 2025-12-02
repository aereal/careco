import {
  fragmentTotalDistance,
  TotalDrivingDistance,
} from '@/components/TotalDrivingDistance';
import { FragmentType, getFragmentData } from '@/graphql';
import { numberOf } from '@/month';
import { FC } from 'react';
import { DriveRecordChart, fragmentChartDataSeries } from '../DriveRecordChart';
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
    <div>
      <h1 className='font-bold text-2xl -mb-4'>
        {year}年{numberOf(month)}月の走行記録
      </h1>
      <TotalDrivingDistance {...summary} />
      <div className='my-8 h-[70vh]'>
        <DriveRecordChart records={records} />
      </div>
    </div>
  );
};
