import { Result } from '@praha/byethrow';
import { describe, expect, test } from 'vitest';
import { EnvironmentVariableNotDefined, parseEnv, type Config } from './config';

describe('parseEnv', () => {
  test('ok', () => {
    expect(
      parseEnv({
        VITE_AUTH0_CLIENT_ID: 'dummy',
        VITE_AUTH0_DOMAIN: 'auth0.test',
        VITE_AUTH0_AUDIENCE: 'https://app.test',
        VITE_AUTH0_ISSUER: 'https://issuer.auth0.test/',
        VITE_BACKEND_ENDPOINT: 'http://upstream.test/graphql',
      }),
    ).toStrictEqual(
      Result.succeed({
        auth0ClientID: 'dummy',
        auth0Domain: 'auth0.test',
        auth0Audience: 'https://app.test',
        auth0Issuer: 'https://issuer.auth0.test/',
        backendEndpoint: 'http://upstream.test/graphql',
      } satisfies Config),
    );
  });

  test('omits VITE_AUTH0_ISSUER', () => {
    expect(
      parseEnv({
        VITE_AUTH0_CLIENT_ID: 'dummy',
        VITE_AUTH0_DOMAIN: 'auth0.test',
        VITE_AUTH0_AUDIENCE: 'https://app.test',
        VITE_BACKEND_ENDPOINT: 'http://upstream.test/graphql',
      }),
    ).toStrictEqual(
      Result.succeed({
        auth0ClientID: 'dummy',
        auth0Domain: 'auth0.test',
        auth0Audience: 'https://app.test',
        auth0Issuer: 'https://auth0.test/',
        backendEndpoint: 'http://upstream.test/graphql',
      } satisfies Config),
    );
  });

  test('none', () => {
    expect(parseEnv({})).toStrictEqual(
      Result.fail([
        new EnvironmentVariableNotDefined('VITE_AUTH0_CLIENT_ID'),
        new EnvironmentVariableNotDefined('VITE_AUTH0_DOMAIN'),
        new EnvironmentVariableNotDefined('VITE_AUTH0_AUDIENCE'),
        new EnvironmentVariableNotDefined('VITE_BACKEND_ENDPOINT'),
      ]),
    );
  });
});
