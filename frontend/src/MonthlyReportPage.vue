<script lang="ts">
import MonthlyReport from '@/components/MonthlyReport.vue';
import { graphql } from '@/graphql';
import { getFirstParam } from '@/utils/get-first-param';
import { formatMonth } from '@/utils/month';
import { parseDateParams } from '@/utils/parse-date-params';
import type { ResultOf } from '@graphql-typed-document-node/core';
import { Result } from '@praha/byethrow';
import { CombinedError, useQuery } from '@urql/vue';
import { ref } from 'vue';
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
  components: { MonthlyReport },
  setup() {
    const { params } = useRoute('/reports/:year/:month');
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
        });
        return { fetching, data, error };
      }),
      Result.orElse((error) =>
        Result.succeed({
          error: ref(new CombinedError({ networkError: error })),
          fetching: ref(false),
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
      <div v-if="fetching">Loading...</div>
      <div v-else-if="error">Error: {{ error.message }}</div>
      <div v-else-if="data">
        <MonthlyReport :summary="data.monthlyReport" />
      </div>
    </div>
  </div>
</template>
