import { graphql } from '@/graphql';

export const GetRoot = graphql(`
  query GetRoot($first: Int!) {
    totalStatistics {
      distanceKilometers
    }
    recentDrivingRecords(first: $first) {
      distanceKilometers
      recordedAt
    }
  }
`);
