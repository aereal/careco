export interface Auth0Config {
  readonly domain: string;
  readonly clientID: string;
  readonly clientSecret: string;
  readonly audience: string;
}

const envDomain = 'AUTH0_DOMAIN' as const;
const envClientID = 'AUTH0_CLIENT_ID' as const;
const envClientSecret = 'AUTH0_CLIENT_SECRET' as const;
const envAudience = 'AUTH0_AUDIENCE' as const;

type MandatoryAuth0Env =
  | typeof envAudience
  | typeof envDomain
  | typeof envClientID
  | typeof envClientSecret;

export class MissingAuth0EnvError extends Error {
  constructor(envs: Set<MandatoryAuth0Env>) {
    const envArray = Array.from(envs);
    envArray.sort();
    super(`missing Auth0 env: ${envArray.join(', ')}`);
  }
}

export const provideAuth0Config = (
  env: Record<string, string | undefined>,
): Auth0Config => {
  const envs = new Set<MandatoryAuth0Env>();
  const domain = env[envDomain];
  if (domain === undefined) {
    envs.add(envDomain);
  }
  const clientID = env[envClientID];
  if (clientID === undefined) {
    envs.add(envClientID);
  }
  const clientSecret = env[envClientSecret];
  if (clientSecret === undefined) {
    envs.add(envClientSecret);
  }
  const audience = env[envAudience];
  if (audience === undefined) {
    envs.add(envAudience);
  }
  if (envs.size > 0) {
    throw new MissingAuth0EnvError(envs);
  }
  return {
    domain: domain!,
    clientID: clientID!,
    clientSecret: clientSecret!,
    audience: audience!,
  };
};
