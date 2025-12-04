import { M2MTokenProvider } from '@/authz';
import { createClient, fetchExchange } from '@urql/core';
import { authExchange } from '@urql/exchange-auth';
import { registerUrql } from '@urql/next/rsc';

const isServerSide = () => typeof window === 'undefined';

const tokenProvider = M2MTokenProvider.fromEnv();

const makeClient = () =>
  createClient({
    url: 'http://localhost:3000/api/graphql',
    exchanges: [
      authExchange(async (utils) => ({
        addAuthToOperation(op) {
          if (isServerSide()) {
            const token = tokenProvider.getToken();
            if (token !== null) {
              return utils.appendHeaders(op, {
                authorization: `Bearer ${token}`,
              });
            }
          }
          return op;
        },
        willAuthError() {
          return tokenProvider.willExpire();
        },
        didAuthError(error) {
          return error.response?.status === 401;
        },
        async refreshAuth() {
          await tokenProvider.refresh();
        },
      })),
      fetchExchange,
    ],
    preferGetMethod: false,
  });

export const { getClient } = registerUrql(makeClient);
