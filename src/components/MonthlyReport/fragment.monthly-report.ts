import { graphql } from '@/graphql';

export const fragmentMonthlySummary = graphql(`
  fragment MonthlySummary on MonthlyReport {
    year
    month
  }
`);
