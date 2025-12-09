import { Auth0Client } from '@auth0/nextjs-auth0/server';
import { provideAuth0Config } from './env';

const {
  clientID: clientId,
  clientSecret,
  domain,
  secret,
} = provideAuth0Config(process.env);

export const client = new Auth0Client({
  clientId,
  clientSecret,
  domain,
  secret,
});
