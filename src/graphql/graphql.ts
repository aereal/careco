/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = T | null | undefined;
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

export type DailyReport = {
  readonly day: Scalars['Int']['output'];
  readonly distanceKilometers: Scalars['Int']['output'];
  readonly memo?: Maybe<Scalars['String']['output']>;
  readonly month: Month;
  readonly recordedAt: Scalars['DateTime']['output'];
  readonly year: Scalars['Int']['output'];
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

export type MonthlyReport = {
  readonly dailyReports: ReadonlyArray<DailyReport>;
  readonly distanceKilometers: Scalars['Int']['output'];
  readonly month: Month;
  readonly year: Scalars['Int']['output'];
};

export type Mutation = {
  readonly recordDrivingRecord: Scalars['Boolean']['output'];
};

export type MutationRecordDrivingRecordArgs = {
  date: Scalars['DateTime']['input'];
  distanceKilometers: Scalars['Int']['input'];
  memo?: InputMaybe<Scalars['String']['input']>;
};

export type Query = {
  readonly monthlyReport: MonthlyReport;
  readonly recentDrivingRecords: ReadonlyArray<DailyReport>;
  readonly totalStatistics: TotalStatistics;
  readonly yearlyReport: YearlyReport;
};

export type QueryMonthlyReportArgs = {
  month: Month;
  year: Scalars['Int']['input'];
};

export type QueryRecentDrivingRecordsArgs = {
  first: Scalars['Int']['input'];
};

export type QueryYearlyReportArgs = {
  year: Scalars['Int']['input'];
};

export type TotalStatistics = {
  readonly distanceKilometers: Scalars['Int']['output'];
};

export type YearlyReport = {
  readonly distanceKilometers: Scalars['Int']['output'];
  readonly monthlyReports: ReadonlyArray<MonthlyReport>;
  readonly year: Scalars['Int']['output'];
};

export type GetRootQueryVariables = Exact<{
  first: Scalars['Int']['input'];
}>;

export type GetRootQuery = {
  readonly totalStatistics: { readonly distanceKilometers: number };
  readonly recentDrivingRecords: ReadonlyArray<{
    readonly distanceKilometers: number;
    readonly recordedAt: Date;
  }>;
};

export type MonthReportQueryVariables = Exact<{
  year: Scalars['Int']['input'];
  month: Month;
}>;

export type MonthReportQuery = {
  readonly monthlyReport: {
    ' $fragmentRefs'?: { MonthlySummaryFragment: MonthlySummaryFragment };
  };
};

export type MonthlySummaryFragment = {
  readonly distanceKilometers: number;
  readonly year: number;
  readonly month: Month;
} & { ' $fragmentName'?: 'MonthlySummaryFragment' };

export const MonthlySummaryFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'MonthlySummary' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'MonthlyReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'distanceKilometers' },
          },
          { kind: 'Field', name: { kind: 'Name', value: 'year' } },
          { kind: 'Field', name: { kind: 'Name', value: 'month' } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<MonthlySummaryFragment, unknown>;
export const GetRootDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'GetRoot' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'first' },
          },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'Int' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'totalStatistics' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'distanceKilometers' },
                },
              ],
            },
          },
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'recentDrivingRecords' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'first' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'first' },
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'distanceKilometers' },
                },
                { kind: 'Field', name: { kind: 'Name', value: 'recordedAt' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetRootQuery, GetRootQueryVariables>;
export const MonthReportDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'MonthReport' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'year' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'Int' } },
          },
        },
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'month' },
          },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'Month' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'monthlyReport' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'year' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'year' },
                },
              },
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'month' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'month' },
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'MonthlySummary' },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'MonthlySummary' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'MonthlyReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'distanceKilometers' },
          },
          { kind: 'Field', name: { kind: 'Name', value: 'year' } },
          { kind: 'Field', name: { kind: 'Name', value: 'month' } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<MonthReportQuery, MonthReportQueryVariables>;
