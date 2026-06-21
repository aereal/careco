import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import { join } from 'node:path';
import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';
import vueRouter from 'vue-router/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vueRouter({
      root: join(process.cwd(), './frontend'),
      dts: './typed-router.d.ts',
    }),
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
