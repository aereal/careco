import { Result } from '@praha/byethrow';
import { describe, expect, test } from 'vitest';
import { InvalidDateParamsError, parseDateParams } from './parse-date-params';

describe('parseDateParams', () => {
  test('ok', () => {
    const got = parseDateParams({ year: '2024', month: '09' });
    expect(got).toStrictEqual(
      Result.succeed(new Date('2024-09-01T00:00:00+09:00')),
    );
  });

  test('invalid date', () => {
    const got = parseDateParams({ year: 'abc', month: '24' });
    expect(got).toStrictEqual(Result.fail(new InvalidDateParamsError()));
  });
});
