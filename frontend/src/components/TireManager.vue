<script setup lang="ts">
import ErrorAlert from '@/components/ErrorAlert.vue';
import { graphql } from '@/graphql';
import { runMutation } from '@/utils/run-operation';
import { Result } from '@praha/byethrow';
import { useMutation, useQuery } from '@urql/vue';
import { format } from 'date-fns/format';
import { formatISO } from 'date-fns/formatISO';
import { parseISO } from 'date-fns/parseISO';
import { computed, ref } from 'vue';

const queryTires = graphql(`
  query ListTires($includeRetired: Boolean!) {
    tires(includeRetired: $includeRetired) {
      id
      name
      purchasedOn
      retired
    }
  }
`);

const mutationAddTire = graphql(`
  mutation AddTire($name: String!, $purchasedOn: DateTime) {
    addTire(name: $name, purchasedOn: $purchasedOn) {
      id
    }
  }
`);

const mutationRetireTire = graphql(`
  mutation RetireTire($id: ID!) {
    retireTire(id: $id) {
      id
    }
  }
`);

const includeRetired = ref(false);
const { data, fetching, error, executeQuery } = useQuery({
  query: queryTires,
  variables: computed(() => ({ includeRetired: includeRetired.value })),
});
const { executeMutation: addTire } = useMutation(mutationAddTire);
const { executeMutation: retireTire } = useMutation(mutationRetireTire);

const name = ref('');
const purchasedOn = ref('');
const active = ref(true);
const mutationError = ref<string | undefined>();

const refetch = (): void => {
  executeQuery({ requestPolicy: 'network-only' });
};

const formatPurchasedOn = (value: Date | string): string =>
  format(new Date(value), 'yyyy/MM/dd');

const handleSubmit = async (e: Event): Promise<void> => {
  e.preventDefault();
  mutationError.value = undefined;
  active.value = false;

  const result = await runMutation(addTire)({
    name: name.value,
    purchasedOn:
      purchasedOn.value === ''
        ? undefined
        : formatISO(parseISO(purchasedOn.value)),
  });
  if (Result.isSuccess(result)) {
    name.value = '';
    purchasedOn.value = '';
    refetch();
  } else {
    mutationError.value = result.error.message;
  }

  active.value = true;
};

const handleRetire = async (id: string, tireName: string): Promise<void> => {
  if (!window.confirm(`${tireName}を廃棄済みにしますか？`)) {
    return;
  }
  mutationError.value = undefined;
  const result = await runMutation(retireTire)({ id });
  if (Result.isSuccess(result)) {
    refetch();
  } else {
    mutationError.value = result.error.message;
  }
};
</script>

<template>
  <div>
    <ErrorAlert :error="{ message: mutationError }" v-if="mutationError" />
    <form class="my-4 flex flex-col gap-2" @submit="handleSubmit">
      <label class="input w-full"
        >名前
        <input
          type="text"
          name="name"
          class="grow"
          required
          v-model="name"
          :disabled="!active"
        />
      </label>
      <label class="input w-full"
        >購入日
        <input
          type="date"
          name="purchasedOn"
          class="grow"
          v-model="purchasedOn"
          :disabled="!active"
        />
      </label>
      <button type="submit" class="btn btn-primary" :disabled="!active">
        追加する
      </button>
    </form>
    <label class="label my-2">
      <input
        type="checkbox"
        name="includeRetired"
        class="checkbox checkbox-sm"
        v-model="includeRetired"
      />
      廃棄済みも表示
    </label>
    <ErrorAlert :error="error" v-if="error" />
    <div class="w-full loading loading-spinner" v-else-if="fetching" />
    <ul class="list" v-else-if="data">
      <li
        class="list-row items-center"
        data-testid="tire-row"
        v-for="tire in data.tires"
        :key="tire.id"
      >
        <div>
          <div>
            {{ tire.name }}
            <span class="badge badge-sm badge-ghost" v-if="tire.retired"
              >廃棄済み</span
            >
          </div>
          <div class="text-xs opacity-60" v-if="tire.purchasedOn">
            購入日: {{ formatPurchasedOn(tire.purchasedOn) }}
          </div>
        </div>
        <button
          type="button"
          class="btn btn-sm btn-ghost"
          v-if="!tire.retired"
          @click="handleRetire(tire.id, tire.name)"
        >
          廃棄
        </button>
      </li>
    </ul>
  </div>
</template>
