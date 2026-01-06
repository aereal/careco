<script lang="ts">
import ErrorAlert from '@/components/ErrorAlert.vue';
import MonthlyReport from '@/components/MonthlyReport.vue';
import { graphql } from '@/graphql';
import { getFirstParam } from '@/utils/get-first-param';
import { formatMonth } from '@/utils/month';
import { parseDateParams } from '@/utils/parse-date-params';
import { useAuth0 } from '@auth0/auth0-vue';
import type { ResultOf } from '@graphql-typed-document-node/core';
import { Result } from '@praha/byethrow';
import { CombinedError, useQuery } from '@urql/vue';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

const queryMonthReport = graphql(`
  query MonthReport($year: Int!, $month: Month!) {
    monthlyReport(year: $year, month: $month) {
      ...MonthlySummary
      ...TotalDistance
      dailyReports {
        ...ChartDataSeries
      }
    }
  }
`);

export default {
  components: { MonthlyReport, ErrorAlert },
  setup() {
    const auth0 = useAuth0();
    const { params } = useRoute('/reports/[year]/[month]');
    const ret = Result.pipe(
      Result.sequence({
        year: getFirstParam('year', params.year),
        month: getFirstParam('month', params.month),
      }),
      Result.andThen(parseDateParams),
      Result.andThen((date) =>
        Result.sequence({
          year: Result.succeed(date.getFullYear()),
          month: formatMonth(date),
        }),
      ),
      Result.map((input) => {
        const { fetching, data, error } = useQuery({
          query: queryMonthReport,
          variables: { year: input.year, month: input.month },
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
        };
      }),
      Result.orElse((error) =>
        Result.succeed({
          authOngoing: auth0.isLoading,
          isAuthenticated: auth0.isAuthenticated,
          error: ref(new CombinedError({ networkError: error })),
          queryFetching: ref(false),
          data: ref<ResultOf<typeof queryMonthReport> | undefined>(undefined),
        }),
      ),
    );
    if (Result.isFailure(ret)) {
      throw ret.error;
    }
    return ret.value;
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
      <ErrorAlert :error="error" v-else-if="error" />
      <div v-else-if="data">
        <MonthlyReport :summary="data.monthlyReport" />
      </div>
    </div>
  </div>
</template>
