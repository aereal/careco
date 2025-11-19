import {
  fragmentTotalDistance,
  TotalDrivingDistance,
} from '@/components/TotalDrivingDistance';
import { FragmentType, getFragmentData } from '@/graphql';
import { numberOf } from '@/month';
import { FC } from 'react';
import { fragmentMonthlySummary } from './fragment.monthly-report';

type Data = FragmentType<typeof fragmentMonthlySummary> &
  FragmentType<typeof fragmentTotalDistance>;

interface MonthlyReportProps {
  readonly summary: Data;
}

export const MonthlyReport: FC<MonthlyReportProps> = ({ summary }) => {
  const { year, month } = getFragmentData(fragmentMonthlySummary, summary);
  return (
    <div>
      <h1 className='font-bold text-2xl -mb-4'>
        {year}年{numberOf(month)}月の走行記録
      </h1>
      <TotalDrivingDistance {...summary} />
    </div>
  );
};
