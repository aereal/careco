<script setup lang="ts">
import { type FragmentType, getFragmentData, graphql } from '@/graphql';

const fragmentDistance = graphql(`
  fragment TotalDistance on DistanceReport {
    odometerValue
  }
`);

const props = defineProps<{
  totalDistance: FragmentType<typeof fragmentDistance>;
}>();

const { odometerValue } = getFragmentData(
  fragmentDistance,
  props.totalDistance,
);
</script>

<template>
  <div className="stats shadow mt-8">
    <div className="stat pl-4">
      <h2 className="stat-title">総走行距離</h2>
      <div className="stat-value">{{ odometerValue }}km</div>
    </div>
  </div>
</template>
