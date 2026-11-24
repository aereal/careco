import type { DocumentNode } from 'graphql';
import { vi, type Mock } from 'vitest';
import { ref, type Ref } from 'vue';

export interface QueryState {
  data: Ref<unknown>;
  fetching: Ref<boolean>;
  error: Ref<unknown>;
  executeQuery: Mock;
}

export const operationName = (doc: DocumentNode): string | undefined => {
  for (const def of doc.definitions) {
    if (def.kind === 'OperationDefinition') {
      return def.name?.value;
    }
  }
  return undefined;
};

export const createQueryState = (data: unknown): QueryState => ({
  data: ref(data),
  fetching: ref(false),
  error: ref(undefined),
  executeQuery: vi.fn(),
});
