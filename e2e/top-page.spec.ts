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
    const totalDistance = mockGetRootData.totalStatistics.cumulativeDistance;
    await expect(page.getByText(totalDistance.toString())).toBeVisible();
  });

  test('最近の走行記録が表示される', async ({ page }) => {
    await page.goto('/');

    // モックデータの最初の記録の距離が表示されることを確認
    const firstRecord = mockGetRootData.recentDrivingRecords.nodes[0];
    await expect(
      page.getByText(firstRecord.cumulativeDistance.toString()),
    ).toBeVisible();
  });
});
