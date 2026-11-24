import { createQueryState, operationName } from '@/test-utils/urql-mock';
import { flushPromises, mount } from '@vue/test-utils';
import type { DocumentNode } from 'graphql';
import { beforeEach, describe, expect, it, vi, type Mock } from 'vitest';
import { isRef } from 'vue';
import TireManager from './TireManager.vue';

const urql = vi.hoisted(() => ({
  queryState: undefined as ReturnType<typeof createQueryState> | undefined,
  queryArgs: undefined as { variables?: unknown } | undefined,
  mutations: {} as Record<
    string,
    Mock<(vars: Record<string, unknown>) => Promise<unknown>>
  >,
}));

vi.mock('@urql/vue', () => ({
  useQuery: (args: { variables?: unknown }) => {
    urql.queryArgs = args;
    return urql.queryState;
  },
  useMutation: (doc: DocumentNode) => {
    const name = operationName(doc) ?? '';
    return { executeMutation: urql.mutations[name] };
  },
}));

const tires = [
  {
    id: 'tire-1',
    name: 'スタッドレス',
    purchasedOn: '2024-11-01T00:00:00+09:00',
    retired: false,
  },
  { id: 'tire-2', name: '夏タイヤ', purchasedOn: null, retired: false },
];

const variables = (): unknown => {
  const v = urql.queryArgs?.variables;
  return isRef(v) ? v.value : v;
};

describe('TireManager', () => {
  beforeEach(() => {
    urql.queryState = createQueryState({ tires });
    urql.mutations = {
      AddTire: vi
        .fn()
        .mockResolvedValue({ data: { addTire: { id: 'tire-3' } } }),
      RetireTire: vi
        .fn()
        .mockResolvedValue({ data: { retireTire: { id: 'tire-1' } } }),
    };
  });

  it('lists tires with their purchase dates', () => {
    const wrapper = mount(TireManager);
    const rows = wrapper.findAll('[data-testid="tire-row"]');
    expect(rows).toHaveLength(2);
    expect(rows[0]?.text()).toContain('スタッドレス');
    expect(rows[0]?.text()).toContain('2024/11/01');
    expect(rows[1]?.text()).toContain('夏タイヤ');
  });

  it('excludes retired tires by default and includes them when toggled', async () => {
    const wrapper = mount(TireManager);
    expect(variables()).toEqual({ includeRetired: false });

    await wrapper.get('input[name="includeRetired"]').setValue(true);
    expect(variables()).toEqual({ includeRetired: true });
  });

  it('shows a badge for retired tires and hides their retire button', () => {
    urql.queryState = createQueryState({
      tires: [
        { id: 'tire-9', name: '古いタイヤ', purchasedOn: null, retired: true },
      ],
    });
    const wrapper = mount(TireManager);
    const row = wrapper.get('[data-testid="tire-row"]');
    expect(row.text()).toContain('廃棄済み');
    expect(row.find('button').exists()).toBe(false);
  });

  it('adds a tire and refetches the list', async () => {
    const wrapper = mount(TireManager);
    await wrapper.get('input[name="name"]').setValue('新しいタイヤ');
    await wrapper.get('input[name="purchasedOn"]').setValue('2025-04-01');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    const addTire = urql.mutations['AddTire'];
    expect(addTire).toHaveBeenCalledTimes(1);
    const vars = addTire?.mock.calls[0]?.[0];
    expect(vars).toMatchObject({ name: '新しいタイヤ' });
    expect(new Date(String(vars?.['purchasedOn']))).toEqual(
      new Date(2025, 3, 1),
    );
    expect(urql.queryState?.executeQuery).toHaveBeenCalledWith({
      requestPolicy: 'network-only',
    });
    expect(
      (wrapper.get('input[name="name"]').element as HTMLInputElement).value,
    ).toBe('');
  });

  it('adds a tire without a purchase date', async () => {
    const wrapper = mount(TireManager);
    await wrapper.get('input[name="name"]').setValue('夏タイヤ2');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(urql.mutations['AddTire']).toHaveBeenCalledWith({
      name: '夏タイヤ2',
      purchasedOn: undefined,
    });
  });

  it('shows an error when adding fails', async () => {
    urql.mutations['AddTire'] = vi
      .fn<(vars: Record<string, unknown>) => Promise<unknown>>()
      .mockResolvedValue({ error: new Error('boom') });
    const wrapper = mount(TireManager);
    await wrapper.get('input[name="name"]').setValue('x');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(wrapper.get('[role="alert"]').text()).toContain('boom');
  });

  it('retires a tire after confirmation and refetches the list', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true);
    const wrapper = mount(TireManager);
    await wrapper.get('[data-testid="tire-row"] button').trigger('click');
    await flushPromises();

    expect(urql.mutations['RetireTire']).toHaveBeenCalledWith({ id: 'tire-1' });
    expect(urql.queryState?.executeQuery).toHaveBeenCalledWith({
      requestPolicy: 'network-only',
    });
  });

  it('does not retire a tire when the confirmation is cancelled', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false);
    const wrapper = mount(TireManager);
    await wrapper.get('[data-testid="tire-row"] button').trigger('click');
    await flushPromises();

    expect(urql.mutations['RetireTire']).not.toHaveBeenCalled();
  });
});
