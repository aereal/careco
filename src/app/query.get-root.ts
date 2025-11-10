import { graphql } from '@/graphql';

export const GetRoot = graphql(`
  query GetRoot($first: Int!) {
    totalStatistics {
      ...TotalDistance
    }
    recentDrivingRecords(first: $first) {
      ...RecordList
    }
  }
`);
