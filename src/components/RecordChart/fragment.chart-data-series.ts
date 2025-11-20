import { graphql } from '@/graphql';

export const fragmentChartDataSeries = graphql(`
  fragment ChartDataSeries on DrivingRecordsConnection {
    nodes {
      distanceKilometers
      recordedAt
      totalDistanceKilometers
    }
  }
`);
