import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';
import { defineConfig } from 'vitest/config';

const inCI = process.env['CI'] !== '';

const config = defineConfig({
  plugins: [tsconfigPaths(), react()],
  test: {
    environment: 'jsdom',
    coverage: {
      enabled: inCI,
      provider: 'v8',
      reporter: ['clover', 'html'],
      reportsDirectory: './coverage.front',
      include: ['./src/**/*.{ts,tsx}'],
      exclude: ['./src/**/*.test.{ts,tsx}'],
    },
  },
});
export default config;
