import { graphql } from '@/graphql';

export const queryMonthReport = graphql(`
  query MonthReport($year: Int!, $month: Month!) {
    monthlyReport(year: $year, month: $month) {
      ...MonthlySummary
    }
  }
`);
