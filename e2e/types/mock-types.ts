/**
 * E2Eテスト用のモックデータ型定義
 * Fragment Maskingを考慮せず、実際のGraphQLレスポンス構造を表現
 */

export type GetRootMockData = {
  totalStatistics: {
    distanceKilometers: number;
  };
  recentDrivingRecords: {
    nodes: Array<{
      distanceKilometers: number;
      recordedAt: string;
    }>;
  };
};
