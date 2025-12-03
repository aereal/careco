import { createClient, fetchExchange } from '@urql/core';
import { registerUrql } from '@urql/next/rsc';

const makeClient = () =>
  createClient({
    url: 'http://localhost:3000/api/graphql',
    exchanges: [fetchExchange],
    preferGetMethod: false,
  });

export const { getClient } = registerUrql(makeClient);
