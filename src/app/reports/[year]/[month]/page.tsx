import {
  fragmentMonthlySummary,
  MonthlyReport,
} from '@/components/MonthlyReport';
import { getFragmentData } from '@/graphql';
import { client } from '@/lib/auth0';
import { numberOf } from '@/month';
import { Result } from '@praha/byethrow';
import { type Metadata, type Route } from 'next';
import { redirect } from 'next/navigation';
import { FC } from 'react';
import { fetchMonthReport } from './fetch-month-report';

export const generateMetadata = async (
  props: PageProps<'/reports/[year]/[month]'>,
): Promise<Metadata> => {
  const session = await client.getSession();
  if (!session?.user) {
    return {};
  }

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
  const session = await client.getSession();
  if (!session?.user) {
    redirect('/auth/login' as Route);
  }
  const ret = await fetchMonthReport(await params);
  if (Result.isFailure(ret)) {
    return <>Error: {ret.error.message}</>;
  }
  const {
    value: {
      data: { monthlyReport },
    },
  } = ret;
  return (
    <div className='max-w-2xl mx-auto'>
      <div className='p-4'>
        <MonthlyReport
          summary={monthlyReport}
          records={monthlyReport.dailyReports}
        />
      </div>
    </div>
  );
};

export default Page;
