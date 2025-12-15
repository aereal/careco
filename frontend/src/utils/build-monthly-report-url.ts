import { format } from 'date-fns/format';

export const buildMonthlyReportURL = (date: Date): string =>
  `/reports/${format(date, 'yyyy/MM')}`;
