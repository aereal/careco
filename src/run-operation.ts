import { Result } from '@praha/byethrow';
import {
  Client,
  DocumentInput,
  type AnyVariables,
  type OperationResult,
} from '@urql/core';

export type SuccessfulOperationResult<D, V extends AnyVariables> = Omit<
  OperationResult<D, V>,
  'error'
>;

export type SuccessfulReadQueryResult<D, V extends AnyVariables> = Omit<
  SuccessfulOperationResult<D, V>,
  'data'
> & { data: D };

export const readQuery = <D, V extends AnyVariables>(
  client: Client,
  doc: DocumentInput<D, V>,
) => {
  const op = runQuery<D, V>(client, doc);
  return (variables: V) =>
    Result.pipe(op(variables), Result.andThen(assumeData));
};

const runQuery = <D, V extends AnyVariables>(
  client: Client,
  doc: DocumentInput<D, V>,
  catchErr: (err: unknown) => Error = toError,
) =>
  Result.try({
    try: async (variables: V): Promise<SuccessfulOperationResult<D, V>> => {
      const { error, ...rest } = await client.query(doc, variables);
      if (error) {
        throw error;
      }
      return rest;
    },
    catch: catchErr,
  });

export class NullDataError extends Error {
  constructor() {
    super();
  }
  public readonly message = 'data is null';
}

const assumeData = <D, V extends AnyVariables>(
  ret: SuccessfulOperationResult<D, V>,
): Result.Result<SuccessfulReadQueryResult<D, V>, NullDataError> =>
  !ret.data
    ? Result.fail(new NullDataError())
    : Result.succeed({
        operation: ret.operation,
        extensions: ret.extensions,
        hasNext: ret.hasNext,
        stale: ret.stale,
        data: ret.data,
      } satisfies SuccessfulReadQueryResult<D, V>);

const toError = (v: unknown): Error =>
  v instanceof Error ? v : new Error(`${v}`);
