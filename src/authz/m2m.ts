import { Auth0Config, provideAuth0Config } from '@/lib/env';
import { TokenProvider } from './provider';
import { VolatileToken } from './volatile-token';

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

type TokenYielderFunc = () => Promise<Auth0TokenResponse>;

export class M2MTokenProvider implements TokenProvider {
  #cachedValue = new VolatileToken();
  #yielder: TokenYielderFunc;

  static fromEnv(): M2MTokenProvider {
    return new M2MTokenProvider(
      m2mTokenYielder(provideAuth0Config(process.env)),
    );
  }

  private constructor(yielder: TokenYielderFunc) {
    this.#yielder = yielder;
  }

  async refresh(): Promise<void> {
    const { access_token, expires_in } = await this.#yielder();
    this.#cachedValue.update(
      access_token,
      new Date(Date.now() + expires_in * 1000),
    );
  }

  willExpire(): boolean {
    return this.#cachedValue.expiredOn(new Date());
  }

  getToken(): string | null {
    return this.#cachedValue.getToken(new Date());
  }
}

export const m2mTokenYielder =
  (opts: Auth0Config): TokenYielderFunc =>
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
