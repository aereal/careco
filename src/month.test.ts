import { Result } from '@praha/byethrow';
import { describe, expect, test } from 'vitest';
import { formatMonth, InvalidMonthNumberError } from './month';

describe('formatMonthNumber', () => {
  test('ok', () => {
    expect(formatMonth(new Date('2025-11-11'))).toStrictEqual(
      Result.succeed('NOVEMBER'),
    );
  });

  test('invalid number', () => {
    expect(formatMonth(new Date('abc'))).toStrictEqual(
      Result.fail(new InvalidMonthNumberError(NaN)),
    );
  });
});
