import { Result } from '@praha/byethrow';
import { isValid } from 'date-fns/isValid';
import { parse } from 'date-fns/parse';

export type Params = {
  readonly year: string;
  readonly month: string;
};

export class InvalidDateParamsError extends Error {
  public readonly message = 'invalid date params';
}

export const parseDateParams = (
  params: Params,
): Result.Result<Date, InvalidDateParamsError> => {
  const parsed = parse(`${params.year}/${params.month}`, 'yyyy/MM', new Date());
  if (!isValid(parsed)) {
    return Result.fail(new InvalidDateParamsError());
  }
  return Result.succeed(parsed);
};
