import { defineConfig, devices, ReporterDescription } from '@playwright/test';
import process from 'node:process';

const PORT = process.env['PORT'] || 5173;
const baseURL = `http://localhost:${PORT}`;
const inCI = (process.env['CI'] ?? '') !== '';

export default defineConfig({
  testDir: './frontend/e2e',
  timeout: 30 * 1000,
  expect: {
    timeout: 5000,
  },
  forbidOnly: inCI,
  retries: 0,
  workers: inCI ? 1 : undefined,
  reporter: [
    ['html'],
    ...(inCI ? ([['github']] satisfies ReporterDescription[]) : []),
  ],
  use: {
    baseURL,
    trace: 'on-first-retry',
    headless: inCI,
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
      },
    },
  ],

  webServer: {
    command: inCI ? 'pnpm run start' : 'pnpm run dev',
    url: baseURL,
    reuseExistingServer: !inCI,
    timeout: 120 * 1000,
  },
});
