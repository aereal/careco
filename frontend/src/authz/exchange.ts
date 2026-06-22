import { type Auth0VueClient } from '@auth0/auth0-vue';
import { authExchange } from '@urql/exchange-auth';
import { ClientTokenProvider } from './client-token-provider';

export const auth0Exchange = (auth0Client: Auth0VueClient) => {
  const provider = new ClientTokenProvider(auth0Client);
  return authExchange(async (utils) => ({
    addAuthToOperation(op) {
      const token = provider.getToken();
      if (token === null) {
        // do not halt the operation even if token is null,
        // let the server return the authorization error.
        return op;
      }
      return utils.appendHeaders(op, {
        authorization: `Bearer ${token}`,
      });
    },
    willAuthError() {
      return provider.willExpire();
    },
    didAuthError(err) {
      return (
        (err.response?.status === 401 ||
          err.networkError?.message?.includes(
            'Unknown or invalid refresh token',
          )) ??
        false
      );
    },
    async refreshAuth() {
      await provider.refresh();
    },
  }));
};
