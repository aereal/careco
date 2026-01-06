<script lang="ts">
import DriveRecordChart from '@/components/DriveRecordChart.vue';
import ErrorAlert from '@/components/ErrorAlert.vue';
import LoginButton from '@/components/LoginButton.vue';
import LogoutButton from '@/components/LogoutButton.vue';
import RecordDialog from '@/components/RecordDialog.vue';
import SelectMonth from '@/components/SelectMonth.vue';
import TotalDistance from '@/components/TotalDistance.vue';
import { graphql } from '@/graphql';
import { useAuth0 } from '@auth0/auth0-vue';
import { useQuery } from '@urql/vue';
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
    DriveRecordChart,
    RecordDialog,
    SelectMonth,
    TotalDistance,
    LoginButton,
    LogoutButton,
    ErrorAlert,
  },
  setup() {
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
  <div class="max-w-2xl mx-auto">
    <div class="p-4">
      <div class="w-full" v-if="authOngoing || queryFetching">
        <div class="w-full loading loading-spinner loading-xl" />
      </div>
      <div v-else-if="!isAuthenticated">
        <div v-if="!isAuthenticated">
          <LoginButton />
        </div>
      </div>
      <div v-else>
        <ErrorAlert :error="error" v-if="error" />
        <div v-else-if="data">
          <TotalDistance
            :total-distance="data.lastReport"
            granularity="day"
            v-if="data.lastReport"
          />
          <div class="my-4">
            <h2 class="font-bold text-md mb-2">月毎の記録を見る</h2>
            <SelectMonth />
          </div>
          <div class="my-8">
            <h1 class="font-bold text-lg -mb-4">最近の記録</h1>
            <div class="my-8 h-[60vh]">
              <DriveRecordChart :data="data.recentDrivingRecords" />
            </div>
          </div>
          <RecordDialog
            :last-odometer-value="data.lastReport"
            @success="handleRecordSuccess"
          />
          <LogoutButton />
        </div>
      </div>
    </div>
  </div>
</template>
