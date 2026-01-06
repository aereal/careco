<script setup lang="ts">
import { type FragmentType, getFragmentData, graphql } from '@/graphql';
import { format } from 'date-fns/format';

const fragmentDistance = graphql(`
  fragment TotalDistance on DistanceReport {
    odometerValue
    reportDate
  }
`);

const props = defineProps<{
  totalDistance: FragmentType<typeof fragmentDistance>;
  granularity: 'day' | 'month';
}>();

const { odometerValue, reportDate } = getFragmentData(
  fragmentDistance,
  props.totalDistance,
);
const formatDate = (d: Date): string =>
  format(d, props.granularity === 'day' ? 'yyyy-MM-dd' : 'yyyy-MM');
</script>

<template>
  <div class="stats shadow mt-8">
    <div class="stat pl-4">
      <h2 class="stat-title">総走行距離</h2>
      <div class="stat-value">{{ odometerValue }}km</div>
      <div class="stat-desc">{{ formatDate(reportDate) }}</div>
    </div>
  </div>
</template>
