import {
  fragmentMonthlySummary,
  MonthlyReport,
} from '@/components/MonthlyReport';
import { getFragmentData } from '@/graphql';
import { numberOf } from '@/month';
import { Result } from '@praha/byethrow';
import { type Metadata } from 'next';
import { FC } from 'react';
import { fetchMonthReport } from './fetch-month-report';

export const generateMetadata = async (
  props: PageProps<'/reports/[year]/[month]'>,
): Promise<Metadata> => {
  const ret = await fetchMonthReport(await props.params);
  if (Result.isFailure(ret)) {
    throw ret.error;
  }
  const { year, month } = getFragmentData(
    fragmentMonthlySummary,
    ret.value.data.monthlyReport,
  );
  return {
    title: `${year}年${numberOf(month)}月の走行記録`,
  };
};

const Page: FC<PageProps<'/reports/[year]/[month]'>> = async ({ params }) => {
  const ret = await fetchMonthReport(await params);
  if (Result.isFailure(ret)) {
    return <>Error: {ret.error.message}</>;
  }
  const {
    value: {
      data: { monthlyReport },
    },
  } = ret;
  return <MonthlyReport {...monthlyReport} />;
};

export default Page;
