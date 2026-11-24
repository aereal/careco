import { flushPromises, mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { createMemoryHistory, createRouter } from 'vue-router';
import BottomNav from './BottomNav.vue';

const createTestRouter = async (initialPath: string) => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/reports/:year/:month', component: { template: '<div />' } },
      { path: '/tires', component: { template: '<div />' } },
    ],
  });
  await router.push(initialPath);
  return router;
};

describe('BottomNav', () => {
  it('renders links for the home and tire items and a button for the report picker', async () => {
    const router = await createTestRouter('/');
    const wrapper = mount(BottomNav, { global: { plugins: [router] } });
    const links = wrapper.findAll('a');
    expect(links.map((a) => a.text())).toEqual(['ホーム', 'タイヤ']);
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('highlights the tire nav item while viewing the tire list', async () => {
    const router = await createTestRouter('/tires');
    const wrapper = mount(BottomNav, { global: { plugins: [router] } });
    const activeItems = wrapper.findAll('.dock-active');
    expect(activeItems).toHaveLength(1);
    expect(activeItems[0]?.text()).toBe('タイヤ');
  });

  it('highlights the report nav item while viewing any report page', async () => {
    const router = await createTestRouter('/reports/2026/06');
    const wrapper = mount(BottomNav, { global: { plugins: [router] } });
    const activeItems = wrapper.findAll('.dock-active');
    expect(activeItems).toHaveLength(1);
    expect(activeItems[0]?.text()).toBe('月別レポート');
  });

  it('opens a month picker dialog and navigates to the chosen month on selection', async () => {
    const router = await createTestRouter('/');
    const wrapper = mount(BottomNav, { global: { plugins: [router] } });

    expect(wrapper.get('dialog').attributes('open')).toBeUndefined();

    await wrapper.get('button').trigger('click');
    expect(wrapper.get('dialog').attributes('open')).toBeDefined();

    await wrapper.get('input[type="date"]').setValue('2024-03-15');
    await wrapper.get('input[type="date"]').trigger('change');
    await flushPromises();

    expect(router.currentRoute.value.fullPath).toBe('/reports/2024/03');
    expect(wrapper.get('dialog').attributes('open')).toBeUndefined();
  });
});
