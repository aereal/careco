/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
  '\n  fragment ChartDataSeries on DrivingRecordsConnection {\n    nodes {\n      odometerValue\n      recordedAt\n      tripDistance\n    }\n  }\n': typeof types.ChartDataSeriesFragmentDoc;
  '\n  fragment MonthlySummary on MonthlyReport {\n    year\n    month\n    ...TotalDistance\n    dailyReports {\n      ...ChartDataSeries\n    }\n  }\n': typeof types.MonthlySummaryFragmentDoc;
  '\n  mutation RecordDrive(\n    $date: DateTime!\n    $distance: Int!\n    $memo: String\n    $tireID: ID!\n  ) {\n    recordDrivingRecord(\n      date: $date\n      odometerValue: $distance\n      memo: $memo\n      tireID: $tireID\n    )\n  }\n': typeof types.RecordDriveDocument;
  '\n  fragment LastOdometerValue on DistanceReport {\n    odometerValue\n  }\n': typeof types.LastOdometerValueFragmentDoc;
  '\n  fragment TireOption on Tire {\n    id\n    name\n  }\n': typeof types.TireOptionFragmentDoc;
  '\n  fragment LastTire on DailyReport {\n    tire {\n      id\n    }\n  }\n': typeof types.LastTireFragmentDoc;
  '\n  query ListTires($includeRetired: Boolean!) {\n    tires(includeRetired: $includeRetired) {\n      id\n      name\n      purchasedOn\n      retired\n    }\n  }\n': typeof types.ListTiresDocument;
  '\n  mutation AddTire($name: String!, $purchasedOn: DateTime) {\n    addTire(name: $name, purchasedOn: $purchasedOn) {\n      id\n    }\n  }\n': typeof types.AddTireDocument;
  '\n  mutation RetireTire($id: ID!) {\n    retireTire(id: $id) {\n      id\n    }\n  }\n': typeof types.RetireTireDocument;
  '\n  fragment TotalDistance on DistanceReport {\n    odometerValue\n    reportDate\n  }\n': typeof types.TotalDistanceFragmentDoc;
  '\n  query GetRoot($first: Int!) {\n    lastReport {\n      ...TotalDistance\n      ...LastOdometerValue\n      ...LastTire\n    }\n    tires {\n      ...TireOption\n    }\n    recentDrivingRecords(first: $first) {\n      ...ChartDataSeries\n    }\n  }\n': typeof types.GetRootDocument;
  '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n      ...TotalDistance\n      dailyReports {\n        ...ChartDataSeries\n      }\n    }\n  }\n': typeof types.MonthReportDocument;
};
const documents: Documents = {
  '\n  fragment ChartDataSeries on DrivingRecordsConnection {\n    nodes {\n      odometerValue\n      recordedAt\n      tripDistance\n    }\n  }\n':
    types.ChartDataSeriesFragmentDoc,
  '\n  fragment MonthlySummary on MonthlyReport {\n    year\n    month\n    ...TotalDistance\n    dailyReports {\n      ...ChartDataSeries\n    }\n  }\n':
    types.MonthlySummaryFragmentDoc,
  '\n  mutation RecordDrive(\n    $date: DateTime!\n    $distance: Int!\n    $memo: String\n    $tireID: ID!\n  ) {\n    recordDrivingRecord(\n      date: $date\n      odometerValue: $distance\n      memo: $memo\n      tireID: $tireID\n    )\n  }\n':
    types.RecordDriveDocument,
  '\n  fragment LastOdometerValue on DistanceReport {\n    odometerValue\n  }\n':
    types.LastOdometerValueFragmentDoc,
  '\n  fragment TireOption on Tire {\n    id\n    name\n  }\n':
    types.TireOptionFragmentDoc,
  '\n  fragment LastTire on DailyReport {\n    tire {\n      id\n    }\n  }\n':
    types.LastTireFragmentDoc,
  '\n  query ListTires($includeRetired: Boolean!) {\n    tires(includeRetired: $includeRetired) {\n      id\n      name\n      purchasedOn\n      retired\n    }\n  }\n':
    types.ListTiresDocument,
  '\n  mutation AddTire($name: String!, $purchasedOn: DateTime) {\n    addTire(name: $name, purchasedOn: $purchasedOn) {\n      id\n    }\n  }\n':
    types.AddTireDocument,
  '\n  mutation RetireTire($id: ID!) {\n    retireTire(id: $id) {\n      id\n    }\n  }\n':
    types.RetireTireDocument,
  '\n  fragment TotalDistance on DistanceReport {\n    odometerValue\n    reportDate\n  }\n':
    types.TotalDistanceFragmentDoc,
  '\n  query GetRoot($first: Int!) {\n    lastReport {\n      ...TotalDistance\n      ...LastOdometerValue\n      ...LastTire\n    }\n    tires {\n      ...TireOption\n    }\n    recentDrivingRecords(first: $first) {\n      ...ChartDataSeries\n    }\n  }\n':
    types.GetRootDocument,
  '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n      ...TotalDistance\n      dailyReports {\n        ...ChartDataSeries\n      }\n    }\n  }\n':
    types.MonthReportDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment ChartDataSeries on DrivingRecordsConnection {\n    nodes {\n      odometerValue\n      recordedAt\n      tripDistance\n    }\n  }\n',
): (typeof documents)['\n  fragment ChartDataSeries on DrivingRecordsConnection {\n    nodes {\n      odometerValue\n      recordedAt\n      tripDistance\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment MonthlySummary on MonthlyReport {\n    year\n    month\n    ...TotalDistance\n    dailyReports {\n      ...ChartDataSeries\n    }\n  }\n',
): (typeof documents)['\n  fragment MonthlySummary on MonthlyReport {\n    year\n    month\n    ...TotalDistance\n    dailyReports {\n      ...ChartDataSeries\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  mutation RecordDrive(\n    $date: DateTime!\n    $distance: Int!\n    $memo: String\n    $tireID: ID!\n  ) {\n    recordDrivingRecord(\n      date: $date\n      odometerValue: $distance\n      memo: $memo\n      tireID: $tireID\n    )\n  }\n',
): (typeof documents)['\n  mutation RecordDrive(\n    $date: DateTime!\n    $distance: Int!\n    $memo: String\n    $tireID: ID!\n  ) {\n    recordDrivingRecord(\n      date: $date\n      odometerValue: $distance\n      memo: $memo\n      tireID: $tireID\n    )\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment LastOdometerValue on DistanceReport {\n    odometerValue\n  }\n',
): (typeof documents)['\n  fragment LastOdometerValue on DistanceReport {\n    odometerValue\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment TireOption on Tire {\n    id\n    name\n  }\n',
): (typeof documents)['\n  fragment TireOption on Tire {\n    id\n    name\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment LastTire on DailyReport {\n    tire {\n      id\n    }\n  }\n',
): (typeof documents)['\n  fragment LastTire on DailyReport {\n    tire {\n      id\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  query ListTires($includeRetired: Boolean!) {\n    tires(includeRetired: $includeRetired) {\n      id\n      name\n      purchasedOn\n      retired\n    }\n  }\n',
): (typeof documents)['\n  query ListTires($includeRetired: Boolean!) {\n    tires(includeRetired: $includeRetired) {\n      id\n      name\n      purchasedOn\n      retired\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  mutation AddTire($name: String!, $purchasedOn: DateTime) {\n    addTire(name: $name, purchasedOn: $purchasedOn) {\n      id\n    }\n  }\n',
): (typeof documents)['\n  mutation AddTire($name: String!, $purchasedOn: DateTime) {\n    addTire(name: $name, purchasedOn: $purchasedOn) {\n      id\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  mutation RetireTire($id: ID!) {\n    retireTire(id: $id) {\n      id\n    }\n  }\n',
): (typeof documents)['\n  mutation RetireTire($id: ID!) {\n    retireTire(id: $id) {\n      id\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment TotalDistance on DistanceReport {\n    odometerValue\n    reportDate\n  }\n',
): (typeof documents)['\n  fragment TotalDistance on DistanceReport {\n    odometerValue\n    reportDate\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  query GetRoot($first: Int!) {\n    lastReport {\n      ...TotalDistance\n      ...LastOdometerValue\n      ...LastTire\n    }\n    tires {\n      ...TireOption\n    }\n    recentDrivingRecords(first: $first) {\n      ...ChartDataSeries\n    }\n  }\n',
): (typeof documents)['\n  query GetRoot($first: Int!) {\n    lastReport {\n      ...TotalDistance\n      ...LastOdometerValue\n      ...LastTire\n    }\n    tires {\n      ...TireOption\n    }\n    recentDrivingRecords(first: $first) {\n      ...ChartDataSeries\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n      ...TotalDistance\n      dailyReports {\n        ...ChartDataSeries\n      }\n    }\n  }\n',
): (typeof documents)['\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n      ...TotalDistance\n      dailyReports {\n        ...ChartDataSeries\n      }\n    }\n  }\n'];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> =
  TDocumentNode extends DocumentNode<infer TType, any> ? TType : never;
