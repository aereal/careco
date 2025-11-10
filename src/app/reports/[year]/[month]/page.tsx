import { MonthlyReport } from '@/components/MonthlyReport';
import { Result } from '@praha/byethrow';
import { FC } from 'react';
import { fetchMonthReport } from './fetch-month-report';

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
