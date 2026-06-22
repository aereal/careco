<script setup lang="ts">
import { buildMonthlyReportURL } from '@/utils/build-monthly-report-url';
import { getEventDateValue } from '@/utils/get-event-date-value';
import { isNavItemActive, navItems } from '@/utils/nav-items';
import { format } from 'date-fns/format';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const items = computed(() =>
  navItems.map((item) => ({
    ...item,
    active: isNavItemActive(route.path, item),
  })),
);

const today = format(new Date(), 'yyyy-MM-dd');
const monthDialogOpen = ref(false);

const openMonthDialog = (e: MouseEvent): void => {
  e.preventDefault();
  monthDialogOpen.value = true;
};
const closeMonthDialog = (e: SubmitEvent): void => {
  e.preventDefault();
  monthDialogOpen.value = false;
};
const handleMonthChange = async (e: Event): Promise<void> => {
  const date = getEventDateValue(e);
  if (date === null) {
    return;
  }
  monthDialogOpen.value = false;
  await router.push(buildMonthlyReportURL(date));
};
</script>

<template>
  <div class="dock">
    <template v-for="item in items" :key="item.label">
      <RouterLink
        v-if="item.type === 'link'"
        :to="item.to"
        :class="{ 'dock-active': item.active }"
      >
        <span class="dock-label">{{ item.label }}</span>
      </RouterLink>
      <button
        v-else
        type="button"
        :class="{ 'dock-active': item.active }"
        @click="openMonthDialog"
      >
        <span class="dock-label">{{ item.label }}</span>
      </button>
    </template>
  </div>
  <dialog class="modal" :open="monthDialogOpen">
    <div class="modal-box">
      <h2 class="font-bold text-lg mb-2">月を選んでレポートを見る</h2>
      <input
        type="date"
        class="input w-full"
        :max="today"
        @change="handleMonthChange"
      />
    </div>
    <form method="dialog" class="modal-backdrop" @submit="closeMonthDialog">
      <button>close month dialog</button>
    </form>
  </dialog>
</template>
