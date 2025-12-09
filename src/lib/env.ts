export interface Auth0Config {
  readonly domain: string;
  readonly clientID: string;
  readonly clientSecret: string;
  readonly audience: string;
  readonly secret: string;
  readonly appBaseURL: string;
}

const envDomain = 'AUTH0_DOMAIN' as const;
const envClientID = 'AUTH0_CLIENT_ID' as const;
const envClientSecret = 'AUTH0_CLIENT_SECRET' as const;
const envAudience = 'AUTH0_AUDIENCE' as const;
const envSecret = 'AUTH0_SECRET' as const;

type MandatoryAuth0Env =
  | typeof envAudience
  | typeof envDomain
  | typeof envClientID
  | typeof envClientSecret
  | typeof envSecret;

export class MissingAuth0EnvError extends Error {
  constructor(envs: Set<MandatoryAuth0Env>) {
    const envArray = Array.from(envs);
    envArray.sort();
    super(`missing Auth0 env: ${envArray.join(', ')}`);
  }
}

const getAppBaseUrl = (env: Record<string, string | undefined>): string => {
  if (env['APP_BASE_URL'] !== undefined) {
    return env['APP_BASE_URL'];
  }
  const url = env['VERCEL_URL'];
  if (url !== undefined) {
    return `https://${url}`;
  }
  throw new Error('no appBaseURL provided');
};

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
  const secret = env[envSecret];
  if (secret === undefined) {
    envs.add(envSecret);
  }
  if (envs.size > 0) {
    throw new MissingAuth0EnvError(envs);
  }
  const appBaseURL = getAppBaseUrl(env);
  return {
    domain: domain!,
    clientID: clientID!,
    clientSecret: clientSecret!,
    audience: audience!,
    secret: secret!,
    appBaseURL,
  };
};
