import { fileURLToPath } from 'node:url';
import { configDefaults, defineConfig, mergeConfig } from 'vitest/config';
import viteConfig from './vite.config';

const inCI = process.env['CI'] !== '';

const testPattern = ['./frontend/src/**/*.test.{ts,tsx}'];

const config = mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'e2e/**'],
      root: fileURLToPath(new URL('./', import.meta.url)),
      coverage: {
        enabled: inCI,
        provider: 'v8',
        reporter: ['clover', 'html'],
        reportsDirectory: './coverage.front',
        include: ['./frontend/src/**/*.{ts,tsx}'],
        exclude: testPattern,
      },
    },
  }),
);

export default config;
