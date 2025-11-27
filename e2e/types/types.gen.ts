import {
  graphql,
  type GraphQLResponseResolver,
  type RequestHandlerOptions,
} from 'msw';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  DateTime: { input: string; output: Date };
};

export type Month =
  | 'APRIL'
  | 'AUGUST'
  | 'DECEMBER'
  | 'FEBRUARY'
  | 'JANUARY'
  | 'JULY'
  | 'JUNE'
  | 'MARCH'
  | 'MAY'
  | 'NOVEMBER'
  | 'OCTOBER'
  | 'SEPTEMBER';

export type GetRootQueryVariables = Exact<{
  first: Scalars['Int']['input'];
}>;

export type GetRootQuery = {
  readonly totalStatistics: { readonly cumulativeDistance: number };
  readonly recentDrivingRecords: {
    readonly nodes: ReadonlyArray<{
      readonly cumulativeDistance: number;
      readonly recordedAt: Date;
    }>;
  };
};

export type MonthReportQueryVariables = Exact<{
  year: Scalars['Int']['input'];
  month: Month;
}>;

export type MonthReportQuery = {
  readonly monthlyReport: {
    readonly year: number;
    readonly month: Month;
    readonly cumulativeDistance: number;
  };
};

type RecordList_DailyReportsConnection_Fragment = {
  readonly nodes: ReadonlyArray<{
    readonly cumulativeDistance: number;
    readonly recordedAt: Date;
  }>;
};

type RecordList_RecentDrivingRecordsConnection_Fragment = {
  readonly nodes: ReadonlyArray<{
    readonly cumulativeDistance: number;
    readonly recordedAt: Date;
  }>;
};

export type RecordListFragment =
  | RecordList_DailyReportsConnection_Fragment
  | RecordList_RecentDrivingRecordsConnection_Fragment;

export type MonthlySummaryFragment = {
  readonly year: number;
  readonly month: Month;
};

export type RecordDriveMutationVariables = Exact<{
  date: Scalars['DateTime']['input'];
  distance: Scalars['Int']['input'];
  memo?: InputMaybe<Scalars['String']['input']>;
}>;

export type RecordDriveMutation = { readonly recordDrivingRecord: boolean };

type TotalDistance_DailyReport_Fragment = {
  readonly cumulativeDistance: number;
};

type TotalDistance_MonthlyReport_Fragment = {
  readonly cumulativeDistance: number;
};

type TotalDistance_TotalStatistics_Fragment = {
  readonly cumulativeDistance: number;
};

type TotalDistance_YearlyReport_Fragment = {
  readonly cumulativeDistance: number;
};

export type TotalDistanceFragment =
  | TotalDistance_DailyReport_Fragment
  | TotalDistance_MonthlyReport_Fragment
  | TotalDistance_TotalStatistics_Fragment
  | TotalDistance_YearlyReport_Fragment;

/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockGetRootQuery(
 *   ({ query, variables }) => {
 *     const { first } = variables;
 *     return HttpResponse.json({
 *       data: { totalStatistics, recentDrivingRecords }
 *     })
 *   },
 *   requestOptions
 * )
 */
export const mockGetRootQuery = (
  resolver: GraphQLResponseResolver<GetRootQuery, GetRootQueryVariables>,
  options?: RequestHandlerOptions,
) =>
  graphql.query<GetRootQuery, GetRootQueryVariables>(
    'GetRoot',
    resolver,
    options,
  );

/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockMonthReportQuery(
 *   ({ query, variables }) => {
 *     const { year, month } = variables;
 *     return HttpResponse.json({
 *       data: { monthlyReport }
 *     })
 *   },
 *   requestOptions
 * )
 */
export const mockMonthReportQuery = (
  resolver: GraphQLResponseResolver<
    MonthReportQuery,
    MonthReportQueryVariables
  >,
  options?: RequestHandlerOptions,
) =>
  graphql.query<MonthReportQuery, MonthReportQueryVariables>(
    'MonthReport',
    resolver,
    options,
  );

/**
 * @param resolver A function that accepts [resolver arguments](https://mswjs.io/docs/api/graphql#resolver-argument) and must always return the instruction on what to do with the intercepted request. ([see more](https://mswjs.io/docs/concepts/response-resolver#resolver-instructions))
 * @param options Options object to customize the behavior of the mock. ([see more](https://mswjs.io/docs/api/graphql#handler-options))
 * @see https://mswjs.io/docs/basics/response-resolver
 * @example
 * mockRecordDriveMutation(
 *   ({ query, variables }) => {
 *     const { date, distance, memo } = variables;
 *     return HttpResponse.json({
 *       data: { recordDrivingRecord }
 *     })
 *   },
 *   requestOptions
 * )
 */
export const mockRecordDriveMutation = (
  resolver: GraphQLResponseResolver<
    RecordDriveMutation,
    RecordDriveMutationVariables
  >,
  options?: RequestHandlerOptions,
) =>
  graphql.mutation<RecordDriveMutation, RecordDriveMutationVariables>(
    'RecordDrive',
    resolver,
    options,
  );
