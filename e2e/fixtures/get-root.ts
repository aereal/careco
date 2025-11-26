import type { GetRootMockData } from '../types/mock-types';

/**
 * GetRootクエリのモックデータ
 */
export const mockGetRootData: GetRootMockData = {
  totalStatistics: {
    distanceKilometers: 12345,
  },
  recentDrivingRecords: {
    nodes: [
      {
        distanceKilometers: 150,
        recordedAt: '2025-11-26T10:00:00Z',
      },
      {
        distanceKilometers: 200,
        recordedAt: '2025-11-25T09:30:00Z',
      },
      {
        distanceKilometers: 175,
        recordedAt: '2025-11-24T14:15:00Z',
      },
    ],
  },
};
