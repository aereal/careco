import { type GetRootQuery } from '../types';

/**
 * GetRootクエリのモックデータ
 */
export const mockGetRootData: GetRootQuery = {
  totalStatistics: {
    cumulativeDistance: 12345,
  },
  recentDrivingRecords: {
    nodes: [
      {
        cumulativeDistance: 150,
        recordedAt: new Date('2025-11-26T10:00:00Z'),
      },
      {
        cumulativeDistance: 200,
        recordedAt: new Date('2025-11-25T09:30:00Z'),
      },
      {
        cumulativeDistance: 175,
        recordedAt: new Date('2025-11-24T14:15:00Z'),
      },
    ],
  },
};
