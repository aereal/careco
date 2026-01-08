<script setup lang="ts">
import { type FragmentType } from '@/graphql';
import Chart from './Chart.vue';
import { fragmentChartDataSeries } from './fragment.chart-data-series';

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();
</script>

<template>
  <Chart :data="props.data" />
</template>
