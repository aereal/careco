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
  <div class="stats shadow mt-8">
    <div class="stat pl-4">
      <h2 class="stat-title">総走行距離</h2>
      <div class="stat-value">{{ odometerValue }}km</div>
    </div>
  </div>
</template>
