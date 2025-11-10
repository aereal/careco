import { FragmentType, getFragmentData } from '@/graphql';
import { numberOf } from '@/month';
import { FC } from 'react';
import { fragmentMonthlySummary } from './fragment.monthly-report';

export const MonthlyReport: FC<FragmentType<typeof fragmentMonthlySummary>> = (
  fragment,
) => {
  const { year, month, distanceKilometers } = getFragmentData(
    fragmentMonthlySummary,
    fragment,
  );
  return (
    <>
      <h1>
        {year}年{numberOf(month)}月
      </h1>
      <div>{distanceKilometers}km</div>
    </>
  );
};
