import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import { join } from 'node:path';
import { fileURLToPath, URL } from 'node:url';
import vueRouter from 'unplugin-vue-router/vite';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vueRouter({ root: join(process.cwd(), './frontend') }),
    vue(),
    vueDevTools(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./frontend/src', import.meta.url)),
    },
  },
});
