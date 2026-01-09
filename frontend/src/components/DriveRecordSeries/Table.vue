<script setup lang="ts">
import { getFragmentData, type FragmentType } from '@/graphql';
import { format } from 'date-fns/format';
import { fragmentChartDataSeries } from './fragment.chart-data-series';

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();
const data = getFragmentData(fragmentChartDataSeries, props.data).nodes;
</script>

<template>
  <table class="table">
    <thead>
      <tr>
        <th>日付</th>
        <th>総走行距離</th>
        <th>増分</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="{ recordedAt, odometerValue, tripDistance } in data"
        :key="recordedAt.valueOf()"
      >
        <td>{{ format(recordedAt, 'yyyy-MM-dd') }}</td>
        <td>{{ odometerValue }}</td>
        <td :class="tripDistance > 0 ? 'text-success-content' : ''">
          {{ tripDistance > 0 ? '+' : '' }}{{ tripDistance }}
        </td>
      </tr>
    </tbody>
  </table>
</template>
