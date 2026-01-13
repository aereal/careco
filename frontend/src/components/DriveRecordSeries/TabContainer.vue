<script setup lang="ts">
import { type FragmentType } from '@/graphql';
import { ref } from 'vue';
import Chart from './Chart.vue';
import { fragmentChartDataSeries } from './fragment.chart-data-series';
import Table from './Table.vue';

const tabChart = 'chart' as const;
const tabTable = 'table' as const;

interface Tab {
  id: typeof tabChart | typeof tabTable;
  label: string;
  content: typeof Table | typeof Chart;
}

const tabs: Array<Tab> = [
  { id: tabChart, label: '推移', content: Chart },
  { id: tabTable, label: 'データ', content: Table },
];

const selectedTab = ref<Tab['id']>(tabChart);
const handleTabChange = (e: Event): void => {
  if (e.target === null) {
    return;
  }
  const { value } = e.target as HTMLInputElement;
  switch (value) {
    case tabTable:
      selectedTab.value = tabTable;
      return;
    case tabChart:
      selectedTab.value = tabChart;
      return;
  }
};

const props = defineProps<{
  data: FragmentType<typeof fragmentChartDataSeries>;
}>();
</script>

<template>
  <div class="tabs tabs-border w-full">
    <template v-for="{ content, id, label } in tabs" :key="id">
      <input
        type="radio"
        name="tabs"
        class="tab"
        :aria-label="label"
        :value="id"
        :checked="selectedTab === id"
        @change="handleTabChange"
      />
      <div class="tab-content border-base-300 bg-base-100 w-full p-4">
        <div class="h-[60vh]">
          <component :is="content" :data="props.data" />
        </div>
      </div>
    </template>
  </div>
</template>
