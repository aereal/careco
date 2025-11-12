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

export const MonthlyReport: FC<Data> = (props) => {
  const { year, month } = getFragmentData(fragmentMonthlySummary, props);
  return (
    <>
      <h1>
        {year}年{numberOf(month)}月
      </h1>
      <TotalDrivingDistance {...props} />
    </>
  );
};
