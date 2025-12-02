'use client';

import { FragmentType, getFragmentData } from '@/graphql';
import { format } from 'date-fns/format';
import { isFirstDayOfMonth } from 'date-fns/isFirstDayOfMonth';
import { type FC } from 'react';
import {
  Area,
  Bar,
  ComposedChart,
  Legend,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';
import { fragmentChartDataSeries } from './fragment.chart-data';
import * as keys from './keys';

interface DriveRecordChartProps {
  readonly records: FragmentType<typeof fragmentChartDataSeries>;
  readonly showMonth?: boolean;
}

export const DriveRecordChart: FC<DriveRecordChartProps> = ({
  records,
  showMonth,
}) => {
  const { nodes } = getFragmentData(fragmentChartDataSeries, records);
  const data = nodes.map(({ recordedAt, odometerValue, tripDistance }) => ({
    [keys.date]: format(recordedAt, 'yyyy-MM-dd'),
    [keys.odometerValue]: odometerValue,
    [keys.tripDistance]: tripDistance,
  }));
  const formatDateTick = (value: Date): string =>
    format(value, showMonth && isFirstDayOfMonth(value) ? 'yyyy-MM-dd' : 'd');
  return (
    <ComposedChart
      data={data}
      className='w-full h-full'
      layout='vertical'
      responsive
    >
      <Legend />
      <Tooltip />
      <YAxis
        type='category'
        dataKey={keys.date}
        tickFormatter={formatDateTick}
        width='auto'
      />
      <XAxis
        type='number'
        width='auto'
        scale='linear'
        unit='km'
        xAxisId={keys.tripDistance}
        orientation='top'
      />
      <XAxis
        type='number'
        width='auto'
        scale='linear'
        unit='km'
        xAxisId={keys.odometerValue}
        orientation='bottom'
      />
      <Area
        type='monotone'
        dataKey={keys.odometerValue}
        xAxisId={keys.odometerValue}
        fill='var(--color-secondary-content)'
        stroke='var(--color-secondary)'
      />
      <Bar
        dataKey={keys.tripDistance}
        type='monotone'
        fill='var(--color-primary)'
        xAxisId={keys.tripDistance}
      />
    </ComposedChart>
  );
};
