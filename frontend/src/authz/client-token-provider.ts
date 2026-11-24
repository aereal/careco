import { getConfig } from '@/config';
import { type Auth0VueClient } from '@auth0/auth0-vue';
import { VolatileToken } from './volatile-token';

export class ClientTokenProvider {
  #cachedValue: VolatileToken = new VolatileToken();
  #client: Auth0VueClient;

  constructor(client: Auth0VueClient) {
    this.#client = client;
  }

  async refresh(): Promise<void> {
    const ret = await this.#client.getAccessTokenSilently({
      detailedResponse: true,
      cacheMode: 'off',
      authorizationParams: {
        audience: getConfig('auth0Audience'),
      },
    });
    if (ret === undefined) {
      return;
    }
    const { expires_in, access_token } = ret;
    const expiresAt = new Date(Date.now() + expires_in * 1000);
    this.#cachedValue.update(access_token, expiresAt);
  }

  willExpire(): boolean {
    return this.#cachedValue.expiredOn(new Date());
  }

  getToken(): string | null {
    return this.#cachedValue.getToken(new Date());
  }
}
