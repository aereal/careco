import { MonthlyReport } from '@/components/MonthlyReport';
import { getClient } from '@/get-client';
import { formatMonth } from '@/month';
import { parseDateParams } from '@/parse-date-params';
import { readQuery } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { FC } from 'react';
import { queryMonthReport } from './query.month-report';

const Page: FC<PageProps<'/reports/[year]/[month]'>> = async ({ params }) => {
  const { year, month } = await params;
  const doQuery = readQuery(getClient(), queryMonthReport);
  const ret = await Result.pipe(
    parseDateParams({ year, month }),
    Result.andThen((date) =>
      Result.pipe(
        formatMonth(date),
        Result.map((month) => ({ date, month })),
      ),
    ),
    Result.andThen(({ date, month }) =>
      doQuery({ year: date.getFullYear(), month }),
    ),
  );
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
