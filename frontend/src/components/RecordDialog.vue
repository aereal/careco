<script setup lang="ts">
import { getFragmentData, graphql, type FragmentType } from '@/graphql';
import { getEventDateValue } from '@/utils/get-event-date-value';
import { getEventNumberValue } from '@/utils/get-event-number-value';
import { runMutation } from '@/utils/run-operation';
import { Result } from '@praha/byethrow';
import { useMutation } from '@urql/vue';
import { format } from 'date-fns/format';
import { formatISO } from 'date-fns/formatISO';
import { computed, ref } from 'vue';

const mutationRecordDrive = graphql(`
  mutation RecordDrive(
    $date: DateTime!
    $distance: Int!
    $memo: String
    $tireID: ID!
  ) {
    recordDrivingRecord(
      date: $date
      odometerValue: $distance
      memo: $memo
      tireID: $tireID
    )
  }
`);

const fragmentLastOdometerValue = graphql(`
  fragment LastOdometerValue on DistanceReport {
    odometerValue
  }
`);

const fragmentTireOption = graphql(`
  fragment TireOption on Tire {
    id
    name
  }
`);

const fragmentLastTire = graphql(`
  fragment LastTire on DailyReport {
    tire {
      id
    }
  }
`);

const { executeMutation: recordDrive } = useMutation(mutationRecordDrive);

const props = defineProps<{
  lastOdometerValue:
    | FragmentType<typeof fragmentLastOdometerValue>
    | undefined
    | null;
  tires: readonly FragmentType<typeof fragmentTireOption>[];
  lastTire: FragmentType<typeof fragmentLastTire> | undefined | null;
}>();

const { odometerValue: lastOdometerValue } = props.lastOdometerValue
  ? getFragmentData(fragmentLastOdometerValue, props.lastOdometerValue)
  : { odometerValue: 0 };

const tires = props.tires.map((t) => getFragmentData(fragmentTireOption, t));
const lastTireID = props.lastTire
  ? getFragmentData(fragmentLastTire, props.lastTire).tire?.id
  : undefined;

const error = ref<string | undefined>();
const modalOpen = ref(false);
const distance = ref(lastOdometerValue);
const date = ref(new Date());
const memo = ref('');
const active = ref(true);
const tireID = ref(
  tires.some((t) => t.id === lastTireID) ? (lastTireID ?? '') : '',
);
const canSubmit = computed(() => active.value && tireID.value !== '');

const emit = defineEmits<{ success: [] }>();

const handleOpenClick = (e: MouseEvent): void => {
  e.preventDefault();
  modalOpen.value = !modalOpen.value;
};
const handleErrorMessageClose = (e: MouseEvent): void => {
  e.preventDefault();
  error.value = undefined;
};
const handleDistanceChange = (e: Event): void => {
  const num = getEventNumberValue(e);
  if (num === undefined) {
    return;
  }
  distance.value = num;
};
const handleDateChange = (e: Event): void => {
  const dv = getEventDateValue(e);
  if (dv === null) {
    return;
  }
  date.value = dv;
};
const handleMemoChange = (e: Event): void => {
  if (!(e.target instanceof HTMLTextAreaElement)) {
    return;
  }
  memo.value = e.target.value;
};
const handleTireChange = (e: Event): void => {
  if (!(e.target instanceof HTMLSelectElement)) {
    return;
  }
  tireID.value = e.target.value;
};
const handleSubmit = async (e: SubmitEvent): Promise<void> => {
  e.preventDefault();
  error.value = undefined;
  active.value = false;

  const result = await runMutation(recordDrive)({
    date: formatISO(date.value),
    distance: distance.value,
    memo: memo.value === '' ? undefined : memo.value,
    tireID: tireID.value,
  });
  const isSuccess =
    Result.isSuccess(result) && result.value.data
      ? result.value.data.recordDrivingRecord
      : false;
  if (isSuccess) {
    modalOpen.value = false;
    date.value = new Date();
    memo.value = '';
    emit('success');
  }
  if (Result.isFailure(result)) {
    error.value = result.error.message;
  }

  active.value = true;
};
</script>

<template>
  <div class="fab fab-above-dock">
    <button class="btn btn-lg btn-circle btn-primary" @click="handleOpenClick">
      +
    </button>
  </div>
  <dialog class="modal" :open="modalOpen">
    <div class="modal-box">
      <form @submit="handleSubmit">
        <label class="input mt-1"
          >走行距離
          <input
            type="number"
            class="grow text-right"
            @change="handleDistanceChange"
            :value="distance"
            :disabled="!active"
          />
          <span class="label">km</span>
        </label>
        <label class="input my-1"
          >日時
          <input
            type="date"
            @change="handleDateChange"
            :value="format(date, 'yyyy-MM-dd')"
            :disabled="!active"
        /></label>
        <label class="select my-1" v-if="tires.length > 0"
          >タイヤ
          <select
            name="tire"
            required
            @change="handleTireChange"
            :value="tireID"
            :disabled="!active"
          >
            <option value="" disabled>選択してください</option>
            <option v-for="tire in tires" :key="tire.id" :value="tire.id">
              {{ tire.name }}
            </option>
          </select>
        </label>
        <p class="my-1 text-sm" v-else>
          <RouterLink to="/tires" class="link">タイヤを登録</RouterLink
          >してから記録してください
        </p>
        <details class="collapse collapse-arrow my-1 p-1 bg-base-100 border">
          <summary class="collapse-title text-sm pl-2 py-2">メモ</summary>
          <div class="collapse-content px-2">
            <textarea
              class="textarea h-24"
              @change="handleMemoChange"
              :value="memo"
              :disabled="!active"
            />
          </div>
        </details>
        <div class="modal-action">
          <button type="submit" class="btn btn-primary" :disabled="!canSubmit">
            記録する
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close record dialog</button>
    </form>
  </dialog>
  <div class="toast toast-above-dock" v-if="error">
    <div class="alert alert-error" role="alert">
      <span>{{ error }}</span>
      <div>
        <button class="btn btn-xs" @click="handleErrorMessageClose">x</button>
      </div>
    </div>
  </div>
</template>
