import { expect, test } from '@playwright/test';

test.describe('ログイン画面', () => {
  test('Loginボタンが表示される', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('button', { name: 'Login' })).toBeVisible();
  });
});
