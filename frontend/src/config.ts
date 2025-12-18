import { Result } from '@praha/byethrow';

export interface Config {
  readonly auth0ClientID: string;
  readonly auth0Domain: string;
  readonly auth0Audience: string;
  readonly auth0Issuer: string;
  readonly backendEndpoint: string;
}

type Environment = Record<string, unknown>;

export class EnvironmentVariableNotDefined extends Error {
  constructor(name: string) {
    super(`environment variable (${name}) is not defined`);
  }
}

const getEnv = (
  env: Environment,
  name: string,
): Result.Result<unknown, EnvironmentVariableNotDefined> =>
  name in env
    ? Result.succeed(env[name])
    : Result.fail(new EnvironmentVariableNotDefined(name));

const castAsString = (value: unknown): Result.Result<string, TypeError> =>
  typeof value === 'string'
    ? Result.succeed(value)
    : Result.fail(new TypeError(`expected string but got ${typeof value}`));

const getString = (env: Environment, name: string) =>
  Result.pipe(getEnv(env, name), Result.andThen(castAsString));

export const parseEnv = (
  env: Environment,
): Result.Result<
  Record<ConfigKey, string>,
  (EnvironmentVariableNotDefined | TypeError)[]
> =>
  Result.pipe(
    Result.collect({
      auth0ClientID: getString(env, 'VITE_AUTH0_CLIENT_ID'),
      auth0Domain: getString(env, 'VITE_AUTH0_DOMAIN'),
      auth0Audience: getString(env, 'VITE_AUTH0_AUDIENCE'),
      backendEndpoint: getString(env, 'VITE_BACKEND_ENDPOINT'),
    } satisfies Record<
      Exclude<ConfigKey, 'auth0Issuer'>,
      ReturnType<typeof getString>
    >),
    Result.andThen((partial) =>
      Result.pipe(
        getString(env, 'VITE_AUTH0_ISSUER'),
        Result.orElse(() => Result.succeed(`https://${partial.auth0Domain}/`)),
        Result.map(
          (issuer) => ({ ...partial, auth0Issuer: issuer }) satisfies Config,
        ),
      ),
    ),
  );

type ConfigKey = keyof Config;

/**
 * @throws {EnvironmentVariableNotDefined|TypeError}
 */
export const getConfig = (key: ConfigKey): string =>
  Result.unwrap(parseEnv(import.meta.env))[key];
