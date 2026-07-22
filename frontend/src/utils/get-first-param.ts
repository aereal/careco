import { Result } from '@praha/byethrow';

export class NoParameterGivenError extends Error {
  public readonly parameterName: string;

  constructor(paramName: string) {
    super(`no parameter given: ${paramName}`);
    this.parameterName = paramName;
  }
}

export const getFirstParam = Result.fn({
  try: (name: string, p: string | string[] | undefined): string => {
    if (p === undefined) {
      throw new NoParameterGivenError(name);
    }
    if (p instanceof Array) {
      if (p.length > 0) {
        return p[0]!;
      } else {
        throw new NoParameterGivenError(name);
      }
    }
    return p;
  },
  catch: (err: unknown): Error =>
    err instanceof Error ? err : new Error(`${err}`),
});
