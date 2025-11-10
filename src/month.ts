import { Month } from '@/graphql/graphql';
import { Result } from '@praha/byethrow';

const months = [
  'JANUARY',
  'FEBRUARY',
  'MARCH',
  'APRIL',
  'MAY',
  'JUNE',
  'JULY',
  'AUGUST',
  'SEPTEMBER',
  'OCTOBER',
  'NOVEMBER',
  'DECEMBER',
] as Array<Month>;

const num2month = new Map<number, Month>(months.map((m, i) => [i + 1, m]));

const month2num = new Map<Month, number>(months.map((m, i) => [m, i + 1]));

export class InvalidMonthNumberError extends Error {
  constructor(public readonly monthNumber: number) {
    super();
  }

  get message(): string {
    return `invalid month number: ${this.monthNumber}`;
  }
}

export const formatMonth = (
  date: Date,
): Result.Result<Month, InvalidMonthNumberError> => {
  const monthNum = date.getMonth() + 1;
  const v = num2month.get(monthNum);
  if (!v) {
    return Result.fail(new InvalidMonthNumberError(monthNum));
  }
  return Result.succeed(v);
};

export const numberOf = (month: Month): number => month2num.get(month)!;
