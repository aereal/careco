<script setup lang="ts">
import { getFragmentData, type FragmentType } from '@/graphql';
import { format } from 'date-fns/format';
import { BarChart, LineChart } from 'echarts/charts';
import {
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from 'echarts/components';
import { use } from 'echarts/core';
import { SVGRenderer } from 'echarts/renderers';
import { computed } from 'vue';
import VChart from 'vue-echarts';
import { fragmentChartDataSeries } from './fragment.chart-data-series';

use([
  SVGRenderer,
  BarChart,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
]);

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();

const data = getFragmentData(fragmentChartDataSeries, props.data);
const option = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: [
    { type: 'value', name: '走行距離 (km)', position: 'top' },
    { type: 'value', name: '総走行距離 (km)', position: 'bottom' },
  ],
  yAxis: {
    type: 'category',
    data: data.nodes.map((n) => format(n.recordedAt, 'yyyy-MM-dd')),
  },
  series: [
    {
      name: '走行距離',
      type: 'bar',
      xAxisIndex: 0,
      data: data.nodes.map((n) => n.tripDistance),
    },
    {
      name: '総走行距離',
      type: 'line',
      xAxisIndex: 1,
      areaStyle: {},
      data: data.nodes.map((n) => n.odometerValue),
    },
  ],
}));
</script>

<template>
  <VChart class="w-full h-full" :option="option" />
</template>
