import { client } from '@/lib/auth0';
import { TokenProvider } from './provider';
import { VolatileToken } from './volatile-token';

export class ClientTokenProvider implements TokenProvider {
  #cachedValue: VolatileToken = new VolatileToken();

  async refresh(): Promise<void> {
    const { token, expiresAt: expEpoch } = await client.getAccessToken();
    this.#cachedValue.update(token, new Date(expEpoch));
  }

  willExpire(): boolean {
    return this.#cachedValue.expiredOn(new Date());
  }

  getToken(): string | null {
    return this.#cachedValue.getToken(new Date());
  }
}
