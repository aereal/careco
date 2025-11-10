import { getClient } from '@/get-client';
import { formatMonth } from '@/month';
import { parseDateParams } from '@/parse-date-params';
import { readQuery } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { queryMonthReport } from './query.month-report';

export interface Params {
  readonly year: string;
  readonly month: string;
}

export const fetchMonthReport = (params: Params) => {
  const { year, month } = params;
  const doQuery = readQuery(getClient(), queryMonthReport);
  return Result.pipe(
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
};
