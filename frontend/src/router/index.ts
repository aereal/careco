import MonthlyReportPage from '@/MonthlyReportPage.vue';
import RootPage from '@/RootPage.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: RootPage,
    },
    {
      path: '/reports/:year/:month',
      component: MonthlyReportPage,
    },
  ],
});

export default router;
