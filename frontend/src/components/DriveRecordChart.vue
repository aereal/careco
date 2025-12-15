<script setup lang="ts">
import { getFragmentData, graphql, type FragmentType } from '@/graphql';

const fragmentChartDataSeries = graphql(`
  fragment ChartDataSeries on DrivingRecordsConnection {
    nodes {
      odometerValue
      recordedAt
      tripDistance
    }
  }
`);

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();

const data = getFragmentData(fragmentChartDataSeries, props.data);
</script>

<template>
  <div>{{ data.nodes.length }} data</div>
</template>
