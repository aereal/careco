import { graphql } from '@/graphql';

export const fragmentRecordList = graphql(`
  fragment RecordList on DrivingRecordsConnection {
    nodes {
      cumulativeDistance
      recordedAt
    }
  }
`);
