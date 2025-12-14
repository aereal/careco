import urql, { fetchExchange } from '@urql/vue';
import { createApp } from 'vue';
import App from './App.vue';
import router from './router';

const app = createApp(App);

app.use(router);
app.use(urql, {
  url: 'http://localhost:8080/graphql',
  exchanges: [fetchExchange],
  preferGetMethod: false,
});

app.mount('#app');
