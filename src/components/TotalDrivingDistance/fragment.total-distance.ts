import { graphql } from '@/graphql';

export const fragmentTotalDistance = graphql(`
  fragment TotalDistance on TotalStatistics {
    distanceKilometers
  }
`);
