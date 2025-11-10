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
  '\n  query GetRoot($first: Int!) {\n    totalStatistics {\n      distanceKilometers\n    }\n    recentDrivingRecords(first: $first) {\n      distanceKilometers\n      recordedAt\n    }\n  }\n': typeof types.GetRootDocument;
  '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n    }\n  }\n': typeof types.MonthReportDocument;
  '\n  fragment MonthlySummary on MonthlyReport {\n    distanceKilometers\n    year\n    month\n  }\n': typeof types.MonthlySummaryFragmentDoc;
};
const documents: Documents = {
  '\n  query GetRoot($first: Int!) {\n    totalStatistics {\n      distanceKilometers\n    }\n    recentDrivingRecords(first: $first) {\n      distanceKilometers\n      recordedAt\n    }\n  }\n':
    types.GetRootDocument,
  '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n    }\n  }\n':
    types.MonthReportDocument,
  '\n  fragment MonthlySummary on MonthlyReport {\n    distanceKilometers\n    year\n    month\n  }\n':
    types.MonthlySummaryFragmentDoc,
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
  source: '\n  query GetRoot($first: Int!) {\n    totalStatistics {\n      distanceKilometers\n    }\n    recentDrivingRecords(first: $first) {\n      distanceKilometers\n      recordedAt\n    }\n  }\n',
): (typeof documents)['\n  query GetRoot($first: Int!) {\n    totalStatistics {\n      distanceKilometers\n    }\n    recentDrivingRecords(first: $first) {\n      distanceKilometers\n      recordedAt\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n    }\n  }\n',
): (typeof documents)['\n  query MonthReport($year: Int!, $month: Month!) {\n    monthlyReport(year: $year, month: $month) {\n      ...MonthlySummary\n    }\n  }\n'];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(
  source: '\n  fragment MonthlySummary on MonthlyReport {\n    distanceKilometers\n    year\n    month\n  }\n',
): (typeof documents)['\n  fragment MonthlySummary on MonthlyReport {\n    distanceKilometers\n    year\n    month\n  }\n'];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> =
  TDocumentNode extends DocumentNode<infer TType, any> ? TType : never;
