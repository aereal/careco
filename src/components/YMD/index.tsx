import { format } from 'date-fns/format';
import { startOfDay } from 'date-fns/startOfDay';
import { FC } from 'react';

export interface YMDProps {
  readonly date: Date;
}

export const YMD: FC<YMDProps> = ({ date }) => {
  const formatted = format(startOfDay(date), 'yyyy-MM-dd');
  return <time dateTime={formatted}>{formatted}</time>;
};
