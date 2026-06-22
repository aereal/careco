<script lang="ts">
import { DriveRecordSeries } from '@/components/DriveRecordSeries';
import ErrorAlert from '@/components/ErrorAlert.vue';
import RecordDialog from '@/components/RecordDialog.vue';
import TotalDistance from '@/components/TotalDistance.vue';
import { graphql } from '@/graphql';
import { useAuth0 } from '@auth0/auth0-vue';
import { useQuery } from '@urql/vue';
import { useTitle } from '@vueuse/core';
import { computed } from 'vue';

const queryGetRoot = graphql(`
  query GetRoot($first: Int!) {
    lastReport {
      ...TotalDistance
      ...LastOdometerValue
    }
    recentDrivingRecords(first: $first) {
      ...ChartDataSeries
    }
  }
`);

export default {
  components: {
    RecordDialog,
    TotalDistance,
    ErrorAlert,
    DriveRecordSeries,
  },
  setup() {
    useTitle('Careco');
    const auth0 = useAuth0();

    const { fetching, data, error, executeQuery } = useQuery({
      query: queryGetRoot,
      variables: { first: 50 },
      pause: computed(
        () => auth0.isLoading.value || !auth0.isAuthenticated.value,
      ),
    });

    return {
      authOngoing: auth0.isLoading,
      isAuthenticated: auth0.isAuthenticated,
      queryFetching: fetching,
      data,
      error,
      handleRecordSuccess(): void {
        executeQuery({ requestPolicy: 'network-only' });
      },
    };
  },
};
</script>

<template>
  <div class="w-full" v-if="authOngoing || queryFetching">
    <div class="w-full loading loading-spinner loading-xl" />
  </div>
  <div v-else-if="!isAuthenticated">ログインしてください</div>
  <div v-else>
    <ErrorAlert :error="error" v-if="error" />
    <div v-else-if="data">
      <TotalDistance
        :total-distance="data.lastReport"
        granularity="day"
        v-if="data.lastReport"
      />
      <div class="my-8">
        <h1 class="font-bold text-lg -mb-4">最近の記録</h1>
        <div class="my-8 h-[60vh]">
          <DriveRecordSeries :data="data.recentDrivingRecords" />
        </div>
      </div>
      <RecordDialog
        :last-odometer-value="data.lastReport"
        @success="handleRecordSuccess"
      />
    </div>
  </div>
</template>
