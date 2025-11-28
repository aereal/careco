import { type GetRootQuery } from '../types';

/**
 * GetRootクエリのモックデータ
 */
export const mockGetRootData: GetRootQuery = {
  totalStatistics: {
    odometerValue: 12345,
  },
  recentDrivingRecords: {
    nodes: [
      {
        odometerValue: 150,
        recordedAt: new Date('2025-11-26T10:00:00Z'),
      },
      {
        odometerValue: 200,
        recordedAt: new Date('2025-11-25T09:30:00Z'),
      },
      {
        odometerValue: 175,
        recordedAt: new Date('2025-11-24T14:15:00Z'),
      },
    ],
  },
};
