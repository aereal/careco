import {
  defineConfig,
  devices,
  type Project,
  type ReporterDescription,
} from '@playwright/test';

const PORT = process.env.PORT || 3000;
const baseURL = `http://localhost:${PORT}`;
const inCI = (process.env['CI'] ?? '') !== '';

const setupProject = {
  name: 'setup',
  testMatch: /.*\.setup\.ts/,
} satisfies Project;

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
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
  },

  projects: [
    setupProject,
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        storageState: 'playwright/.auth/user.json',
      },
      dependencies: [setupProject.name],
    },
  ],

  webServer: {
    command: 'pnpm build && pnpm start',
    url: baseURL,
    reuseExistingServer: !inCI,
    timeout: 120 * 1000,
  },
});
