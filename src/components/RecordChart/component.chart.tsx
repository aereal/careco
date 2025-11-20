'use client';

import { getFragmentData, type FragmentType } from '@/graphql';
import { format } from 'date-fns/format';
import { type FC } from 'react';
import { Bar, BarChart, Legend, Tooltip, XAxis, YAxis } from 'recharts';
import { TooltipContent } from './component.tooltip-content';
import { fragmentChartDataSeries } from './fragment.chart-data-series';
import * as keys from './keys';

interface RecordChartProps {
  readonly records: FragmentType<typeof fragmentChartDataSeries>;
  readonly omitMonth?: boolean;
}

export const RecordChart: FC<RecordChartProps> = ({ records, omitMonth }) => {
  const { nodes } = getFragmentData(fragmentChartDataSeries, records);
  const data = nodes.map((n) => ({
    [keys.date]: format(n.recordedAt, 'yyyy-MM-dd'),
    [keys.distance]: n.distanceKilometers,
    [keys.totalDistance]: n.totalDistanceKilometers,
  }));
  const template = omitMonth ? 'd' : 'yyyy-MM-dd';
  const formatXAxisTick = (value: Date): string => format(value, template);
  return (
    <BarChart
      className='w-full h-full'
      data={data}
      layout='vertical'
      responsive
    >
      <Tooltip content={TooltipContent} />
      <Legend />
      <YAxis
        dataKey={keys.date}
        tickFormatter={formatXAxisTick}
        type='category'
      />
      <XAxis width='auto' scale='linear' unit='km' type='number' />
      <Bar
        type='monotone'
        dataKey={keys.distance}
        fill='var(--color-primary)'
      />
    </BarChart>
  );
};
