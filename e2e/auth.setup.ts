import { generateSessionCookie } from '@auth0/nextjs-auth0/testing';
import { test } from '@playwright/test';

const authFile = 'playwright/.auth/user.json';

test('authenticate', async ({ page, baseURL }) => {
  const cookie = await generateSessionCookie(
    {
      user: {
        sub: 'test-user-id',
        email: 'test@example.com',
        name: 'Test User',
      },
    },
    {
      secret: process.env.AUTH0_SECRET!,
    },
  );
  await page.context().addCookies([
    {
      name: 'appSession',
      value: cookie,
      domain: 'localhost',
      path: '/',
      httpOnly: true,
      sameSite: 'Lax',
      expires: Math.floor(Date.now() / 1000) + 86400,
    },
  ]);
  await page.goto(baseURL || 'http://localhost:3000');
  await page.context().storageState({ path: authFile });
});
