<script setup lang="ts">
import { buildMonthlyReportURL } from '@/utils/build-monthly-report-url';
import { getEventDateValue } from '@/utils/get-event-date-value';
import { format } from 'date-fns/format';
import { useRouter } from 'vue-router';

const router = useRouter();
const handleChange = async (e: Event): Promise<void> => {
  e.preventDefault();
  const date = getEventDateValue(e);
  if (date === null) {
    return;
  }
  const nextURL = buildMonthlyReportURL(date);
  await router.push(nextURL);
};
const currentMonth = format(new Date(), 'yyyy-MM');
</script>

<template>
  <input type="month" :max="currentMonth" @change="handleChange" />
</template>
