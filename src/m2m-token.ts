class VolatileToken {
  public readonly expiresAt: Date;

  #token: string;

  constructor(token: string, expiresAt: Date) {
    this.#token = token;
    this.expiresAt = expiresAt;
  }

  getToken(date: Date): string | null {
    if (this.expiredOn(date)) {
      return null;
    }
    return this.#token;
  }

  expiredOn(date: Date): boolean {
    return this.expiresAt > date;
  }
}

class ObtainM2MTokenError extends Error {
  public readonly statusCode: number;
  public readonly upstreamError: string;

  constructor(status: number, upstreamError: string) {
    super(`failed to obtain m2m token: ${status}: ${upstreamError}`);
    this.statusCode = status;
    this.upstreamError = upstreamError;
  }
}

interface Auth0TokenResponse {
  readonly access_token: string;
  readonly expires_in: number;
}

interface FetchM2MTokenOptions {
  readonly domain: string;
  readonly clientID: string;
  readonly clientSecret: string;
  readonly audience: string;
}

const envDomain = 'AUTH0_DOMAIN' as const;
const envClientID = 'AUTH0_CLIENT_ID' as const;
const envClientSecret = 'AUTH0_CLIENT_SECRET' as const;
const envAudience = 'AUTH0_AUDIENCE' as const;

type MandantoryAuth0Env =
  | typeof envAudience
  | typeof envDomain
  | typeof envClientID
  | typeof envClientSecret;

class MissingAuth0EnvError extends Error {
  constructor(envs: Set<MandantoryAuth0Env>) {
    const envArray = Array.from(envs);
    envArray.sort();
    super(`missing Auth0 env: ${envArray.join(', ')}`);
  }
}

export const fetchM2MTokenOptionsFromEnv = (
  env: Record<string, string | undefined>,
): FetchM2MTokenOptions => {
  const envs = new Set<MandantoryAuth0Env>();
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

type TokenYielderFunc = () => Promise<Auth0TokenResponse>;

export class M2MTokenProvider {
  #cachedValue: VolatileToken | null = null;
  #yielder: TokenYielderFunc;

  constructor(yielder: TokenYielderFunc) {
    this.#yielder = yielder;
  }

  async refresh(): Promise<void> {
    const { access_token, expires_in } = await this.#yielder();
    this.#cachedValue = new VolatileToken(access_token, new Date(expires_in));
  }

  willExpire(): boolean {
    if (this.#cachedValue === null) {
      return true;
    }
    return this.#cachedValue.expiredOn(new Date());
  }

  getToken(): string | null {
    if (this.#cachedValue === null) {
      return null;
    }
    return this.#cachedValue.getToken(new Date());
  }
}

export const auth0TokenYielder =
  (opts: FetchM2MTokenOptions): TokenYielderFunc =>
  async () => {
    const tokenURL = `https://${opts.domain}/oauth/token`;
    const resp = await fetch(tokenURL, {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify({
        grant_type: 'client_credentials',
        client_id: opts.clientID,
        client_secret: opts.clientSecret,
        audience: opts.audience,
      }),
    });
    if (!resp.ok) {
      throw new ObtainM2MTokenError(resp.status, await resp.text());
    }
    return (await resp.json()) as Auth0TokenResponse;
  };
