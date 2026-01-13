<script setup lang="ts">
import { DriveRecordSeries } from '@/components/DriveRecordSeries';
import TotalDistance from '@/components/TotalDistance.vue';
import { getFragmentData, graphql, type FragmentType } from '@/graphql';
import { numberOf } from '@/utils/month';

const fragmentMonthlySummary = graphql(`
  fragment MonthlySummary on MonthlyReport {
    year
    month
    ...TotalDistance
    dailyReports {
      ...ChartDataSeries
    }
  }
`);

const props = defineProps<{
  summary: FragmentType<typeof fragmentMonthlySummary>;
}>();
const summary = getFragmentData(fragmentMonthlySummary, props.summary);
</script>

<template>
  <div>
    <h1 class="font-bold text-2xl -mb-4">
      {{ summary.year }}年{{ numberOf(summary.month) }}月の走行記録
    </h1>
    <TotalDistance :total-distance="summary" granularity="month" />
    <div class="my-8 h-[70vh]">
      <DriveRecordSeries :data="summary.dailyReports" />
    </div>
  </div>
</template>
