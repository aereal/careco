import { expect, test } from 'next/experimental/testmode/playwright/msw';
import { mockGetRootData } from './fixtures/get-root';
import { handlers } from './mocks/handlers';

/**
 * Next.js experimental testProxy + MSW を使用したテスト
 * Server Component と Client Component 両方のリクエストをモック可能
 */

// グローバルなMSWハンドラーを設定
test.use({
  mswHandlers: [handlers, { scope: 'test' }],
});

test.describe('トップページ', () => {
  test('基本的な要素が表示される', async ({ page }) => {
    await page.goto('/');

    // 主要な見出しが表示されていることを確認
    await expect(page.getByText('最近の記録')).toBeVisible();
    await expect(page.getByText('月毎の記録を見る')).toBeVisible();
  });

  test('総走行距離が表示される', async ({ page }) => {
    await page.goto('/');

    // モックデータの総走行距離が表示されることを確認
    const totalDistance = mockGetRootData.totalStatistics.odometerValue;
    await expect(page.getByText(totalDistance.toString())).toBeVisible();
  });

  test('最近の走行記録が表示される', async ({ page }) => {
    await page.goto('/');

    // Rechartsのグラフは role="application" を持つSVG要素として描画される
    const chartSvg = page.getByRole('application');
    await expect(chartSvg).toBeVisible();
  });
});
