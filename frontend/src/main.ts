import { auth0Exchange } from '@/authz/exchange';
import { createAuth0 } from '@auth0/auth0-vue';
import urql, { fetchExchange } from '@urql/vue';
import { createApp } from 'vue';
import App from './App.vue';
import { getConfig } from './config';
import router from './router';
import './style.css';

// oxlint-disable-next-line typescript/no-unsafe-argument -- oxlint's type-aware linting doesn't yet resolve .vue module types (oxc-project/oxc#21936), so `App` appears untyped here
const app = createApp(App);

const auth0Plugin = createAuth0({
  domain: getConfig('auth0Domain'),
  clientId: getConfig('auth0ClientID'),
  issuer: getConfig('auth0Issuer'),
  legacySameSiteCookie: false,
  useRefreshTokens: true,
  cacheLocation: 'localstorage',
  authorizationParams: {
    display: 'touch',
    prompt: 'login',
    redirect_uri: window.location.origin,
    audience: getConfig('auth0Audience'),
  },
});

app.use(router);
app.use(urql, {
  url: getConfig('backendEndpoint'),
  exchanges: [auth0Exchange(auth0Plugin), fetchExchange],
  preferGetMethod: false,
});
app.use(auth0Plugin);

app.mount('#app');
