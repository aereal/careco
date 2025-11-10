import { graphql } from '@/graphql';

export const mutationRecordDrive = graphql(`
  mutation RecordDrive($date: DateTime!, $distance: Int!, $memo: String) {
    recordDrivingRecord(date: $date, distanceKilometers: $distance, memo: $memo)
  }
`);
