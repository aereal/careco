import { auth0Exchange } from '@/authz/exchange';
import { createAuth0 } from '@auth0/auth0-vue';
import urql, { fetchExchange } from '@urql/vue';
import { createApp } from 'vue';
import App from './App.vue';
import { config } from './config';
import router from './router';
import './style.css';

const app = createApp(App);

const auth0Plugin = createAuth0({
  domain: config.AUTH0_DOMAIN,
  clientId: config.AUTH0_CLIENT_ID,
  legacySameSiteCookie: false,
  useRefreshTokens: true,
  cacheLocation: 'localstorage',
  authorizationParams: {
    display: 'touch',
    prompt: 'consent',
    redirect_uri: window.location.origin,
    audience: config.AUTH0_AUDIENCE,
  },
});

app.use(router);
app.use(urql, {
  url: 'http://localhost:8080/graphql',
  exchanges: [auth0Exchange(auth0Plugin), fetchExchange],
  preferGetMethod: false,
});
app.use(auth0Plugin);

app.mount('#app');
