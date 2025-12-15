<script lang="ts">
import DriveRecordChart from '@/components/DriveRecordChart.vue';
import RecordDialog from '@/components/RecordDialog.vue';
import SelectMonth from '@/components/SelectMonth.vue';
import TotalDistance from '@/components/TotalDistance.vue';
import { graphql } from '@/graphql';
import { useQuery } from '@urql/vue';

const queryGetRoot = graphql(`
  query GetRoot($first: Int!) {
    totalStatistics {
      ...TotalDistance
    }
    recentDrivingRecords(first: $first) {
      ...ChartDataSeries
    }
  }
`);

export default {
  components: {
    DriveRecordChart,
    RecordDialog,
    SelectMonth,
    TotalDistance,
  },
  setup() {
    const { fetching, data, error } = useQuery({
      query: queryGetRoot,
      variables: { first: 50 },
    });
    return { fetching, data, error };
  },
};
</script>

<template>
  <div className="max-w-2xl mx-auto">
    <div className="p-4">
      <div v-if="fetching">Loading...</div>
      <div v-else-if="error">Error: {{ error.message }}</div>
      <div v-else-if="data">
        <TotalDistance :total-distance="data.totalStatistics" />
        <div className="my-4">
          <h2 className="font-bold text-md mb-2">月毎の記録を見る</h2>
          <SelectMonth />
        </div>
        <div className="my-8">
          <h1 className="font-bold text-lg -mb-4">最近の記録</h1>
          <div className="my-8 h-[60vh]">
            <DriveRecordChart :data="data.recentDrivingRecords" />
          </div>
        </div>
        <RecordDialog />
      </div>
    </div>
  </div>
</template>
