<script lang="ts">
import { graphql } from '@/graphql';

export const fragmentChartDataSeries = graphql(`
  fragment ChartDataSeries on DrivingRecordsConnection {
    nodes {
      odometerValue
      recordedAt
      tripDistance
    }
  }
`);
</script>

<script setup lang="ts">
import { getFragmentData, type FragmentType } from '@/graphql';

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();
const data = getFragmentData(fragmentChartDataSeries, props.data);
</script>

<template>
  <div>{{ data.nodes.length }} data</div>
</template>
