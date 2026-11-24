import { makeFragmentData } from '@/graphql';
import {
  LastOdometerValueFragmentDoc,
  LastTireFragmentDoc,
  TireOptionFragmentDoc,
} from '@/graphql/graphql';
import { operationName } from '@/test-utils/urql-mock';
import { flushPromises, mount } from '@vue/test-utils';
import type { DocumentNode } from 'graphql';
import { beforeEach, describe, expect, it, vi, type Mock } from 'vitest';
import { createMemoryHistory, createRouter } from 'vue-router';
import RecordDialog from './RecordDialog.vue';

const urql = vi.hoisted(() => ({
  mutations: {} as Record<
    string,
    Mock<(vars: Record<string, unknown>) => Promise<unknown>>
  >,
}));

vi.mock('@urql/vue', () => ({
  useMutation: (doc: DocumentNode) => ({
    executeMutation: urql.mutations[operationName(doc) ?? ''],
  }),
}));

const tires = [
  { id: 'tire-1', name: 'スタッドレス' },
  { id: 'tire-2', name: '夏タイヤ' },
];

interface DialogInput {
  tires: readonly { id: string; name: string }[];
  lastTireID?: string;
}

const mountDialog = ({ tires: tireList, lastTireID }: DialogInput) => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/tires', component: { template: '<div />' } },
    ],
  });
  return mount(RecordDialog, {
    props: {
      lastOdometerValue: makeFragmentData(
        { odometerValue: 100 },
        LastOdometerValueFragmentDoc,
      ),
      tires: tireList.map((t) => makeFragmentData(t, TireOptionFragmentDoc)),
      lastTire:
        lastTireID === undefined
          ? null
          : makeFragmentData({ tire: { id: lastTireID } }, LastTireFragmentDoc),
    },
    global: { plugins: [router] },
  });
};

const selectedTire = (wrapper: ReturnType<typeof mountDialog>): string =>
  (wrapper.get('select[name="tire"]').element as HTMLSelectElement).value;

describe('RecordDialog', () => {
  beforeEach(() => {
    urql.mutations = {
      RecordDrive: vi
        .fn<(vars: Record<string, unknown>) => Promise<unknown>>()
        .mockResolvedValue({ data: { recordDrivingRecord: true } }),
    };
  });

  it('lists tires as options', () => {
    const wrapper = mountDialog({ tires });
    const options = wrapper
      .findAll('select[name="tire"] option')
      .map((o) => o.text());
    expect(options).toEqual(
      expect.arrayContaining(['スタッドレス', '夏タイヤ']),
    );
  });

  it('preselects the tire used in the last record', () => {
    const wrapper = mountDialog({
      tires,
      lastTireID: 'tire-2',
    });
    expect(selectedTire(wrapper)).toBe('tire-2');
  });

  it('does not preselect a tire that is no longer available', () => {
    const wrapper = mountDialog({
      tires,
      lastTireID: 'tire-retired',
    });
    expect(selectedTire(wrapper)).toBe('');
  });

  it('requires a tire to be selected', async () => {
    const wrapper = mountDialog({ tires });
    expect(
      (wrapper.get('select[name="tire"]').element as HTMLSelectElement)
        .required,
    ).toBe(true);
    expect(
      wrapper.get('button[type="submit"]').attributes('disabled'),
    ).toBeDefined();

    await wrapper.get('select[name="tire"]').setValue('tire-1');
    expect(
      wrapper.get('button[type="submit"]').attributes('disabled'),
    ).toBeUndefined();
  });

  it('sends the selected tire with the record', async () => {
    const wrapper = mountDialog({
      tires,
      lastTireID: 'tire-1',
    });
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(urql.mutations['RecordDrive']).toHaveBeenCalledWith(
      expect.objectContaining({ distance: 100, tireID: 'tire-1' }),
    );
    expect(wrapper.emitted('success')).toHaveLength(1);
  });

  it('links to the tire list when no tires are registered', () => {
    const wrapper = mountDialog({ tires: [] });
    expect(wrapper.find('select[name="tire"]').exists()).toBe(false);
    expect(wrapper.get('a[href="/tires"]').text()).toContain('タイヤを登録');
    expect(
      wrapper.get('button[type="submit"]').attributes('disabled'),
    ).toBeDefined();
  });
});
