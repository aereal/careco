import { Result } from '@praha/byethrow';
import {
  type AnyVariables,
  type OperationContext,
  type OperationResult,
} from '@urql/core';

export type SuccessfulOperationResult<D, V extends AnyVariables> = Omit<
  OperationResult<D, V>,
  'error'
>;

export type OperationFunc<D, V extends AnyVariables> = (
  variables: V,
  context?: Partial<OperationContext>,
) => Promise<OperationResult<D, AnyVariables>>;

export const runMutation = <D, V extends AnyVariables>(
  fn: OperationFunc<D, V>,
  catchErr: (err: unknown) => Error = toError,
): ((
  variables: V,
) => Promise<Result.Result<SuccessfulOperationResult<D, V>, Error>>) =>
  Result.fn({
    try: async (variables: V): Promise<SuccessfulOperationResult<D, V>> => {
      const { error, ...rest } = await fn(variables);
      if (error) {
        throw error;
      }
      return rest as SuccessfulOperationResult<D, V>;
    },
    catch: catchErr,
  });

const toError = (v: unknown): Error =>
  v instanceof Error ? v : new Error(String(v));
