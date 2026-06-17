/* eslint-disable */
/** Internal type. DO NOT USE DIRECTLY. */
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
/** Internal type. DO NOT USE DIRECTLY. */
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never;
    };
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
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

type ChartDataSeries_DailyReportsConnection_Fragment = {
  readonly nodes: ReadonlyArray<{
    readonly odometerValue: number;
    readonly recordedAt: Date;
    readonly tripDistance: number;
  }>;
} & { ' $fragmentName'?: 'ChartDataSeries_DailyReportsConnection_Fragment' };

type ChartDataSeries_RecentDrivingRecordsConnection_Fragment = {
  readonly nodes: ReadonlyArray<{
    readonly odometerValue: number;
    readonly recordedAt: Date;
    readonly tripDistance: number;
  }>;
} & {
  ' $fragmentName'?: 'ChartDataSeries_RecentDrivingRecordsConnection_Fragment';
};

export type ChartDataSeriesFragment =
  | ChartDataSeries_DailyReportsConnection_Fragment
  | ChartDataSeries_RecentDrivingRecordsConnection_Fragment;

export type MonthlySummaryFragment = ({
  readonly year: number;
  readonly month: Month;
  readonly dailyReports: {
    ' $fragmentRefs'?: {
      ChartDataSeries_DailyReportsConnection_Fragment: ChartDataSeries_DailyReportsConnection_Fragment;
    };
  };
} & {
  ' $fragmentRefs'?: {
    TotalDistance_MonthlyReport_Fragment: TotalDistance_MonthlyReport_Fragment;
  };
}) & { ' $fragmentName'?: 'MonthlySummaryFragment' };

export type RecordDriveMutationVariables = Exact<{
  date: string;
  distance: number;
  memo?: string | null | undefined;
}>;

export type RecordDriveMutation = { readonly recordDrivingRecord: boolean };

type LastOdometerValue_DailyReport_Fragment = {
  readonly odometerValue: number;
} & { ' $fragmentName'?: 'LastOdometerValue_DailyReport_Fragment' };

type LastOdometerValue_MonthlyReport_Fragment = {
  readonly odometerValue: number;
} & { ' $fragmentName'?: 'LastOdometerValue_MonthlyReport_Fragment' };

type LastOdometerValue_YearlyReport_Fragment = {
  readonly odometerValue: number;
} & { ' $fragmentName'?: 'LastOdometerValue_YearlyReport_Fragment' };

export type LastOdometerValueFragment =
  | LastOdometerValue_DailyReport_Fragment
  | LastOdometerValue_MonthlyReport_Fragment
  | LastOdometerValue_YearlyReport_Fragment;

type TotalDistance_DailyReport_Fragment = {
  readonly odometerValue: number;
  readonly reportDate: Date;
} & { ' $fragmentName'?: 'TotalDistance_DailyReport_Fragment' };

type TotalDistance_MonthlyReport_Fragment = {
  readonly odometerValue: number;
  readonly reportDate: Date;
} & { ' $fragmentName'?: 'TotalDistance_MonthlyReport_Fragment' };

type TotalDistance_YearlyReport_Fragment = {
  readonly odometerValue: number;
  readonly reportDate: Date;
} & { ' $fragmentName'?: 'TotalDistance_YearlyReport_Fragment' };

export type TotalDistanceFragment =
  | TotalDistance_DailyReport_Fragment
  | TotalDistance_MonthlyReport_Fragment
  | TotalDistance_YearlyReport_Fragment;

export type GetRootQueryVariables = Exact<{
  first: number;
}>;

export type GetRootQuery = {
  readonly lastReport: {
    ' $fragmentRefs'?: {
      TotalDistance_DailyReport_Fragment: TotalDistance_DailyReport_Fragment;
      LastOdometerValue_DailyReport_Fragment: LastOdometerValue_DailyReport_Fragment;
    };
  } | null;
  readonly recentDrivingRecords: {
    ' $fragmentRefs'?: {
      ChartDataSeries_RecentDrivingRecordsConnection_Fragment: ChartDataSeries_RecentDrivingRecordsConnection_Fragment;
    };
  };
};

export type MonthReportQueryVariables = Exact<{
  year: number;
  month: Month;
}>;

export type MonthReportQuery = {
  readonly monthlyReport: {
    readonly dailyReports: {
      ' $fragmentRefs'?: {
        ChartDataSeries_DailyReportsConnection_Fragment: ChartDataSeries_DailyReportsConnection_Fragment;
      };
    };
  } & {
    ' $fragmentRefs'?: {
      MonthlySummaryFragment: MonthlySummaryFragment;
      TotalDistance_MonthlyReport_Fragment: TotalDistance_MonthlyReport_Fragment;
    };
  };
};

export const TotalDistanceFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TotalDistance' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
          { kind: 'Field', name: { kind: 'Name', value: 'reportDate' } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<TotalDistanceFragment, unknown>;
export const ChartDataSeriesFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'ChartDataSeries' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DrivingRecordsConnection' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'nodes' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'odometerValue' },
                },
                { kind: 'Field', name: { kind: 'Name', value: 'recordedAt' } },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'tripDistance' },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<ChartDataSeriesFragment, unknown>;
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
          { kind: 'Field', name: { kind: 'Name', value: 'year' } },
          { kind: 'Field', name: { kind: 'Name', value: 'month' } },
          {
            kind: 'FragmentSpread',
            name: { kind: 'Name', value: 'TotalDistance' },
          },
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'dailyReports' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'ChartDataSeries' },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TotalDistance' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
          { kind: 'Field', name: { kind: 'Name', value: 'reportDate' } },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'ChartDataSeries' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DrivingRecordsConnection' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'nodes' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'odometerValue' },
                },
                { kind: 'Field', name: { kind: 'Name', value: 'recordedAt' } },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'tripDistance' },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<MonthlySummaryFragment, unknown>;
export const LastOdometerValueFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'LastOdometerValue' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<LastOdometerValueFragment, unknown>;
export const RecordDriveDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'RecordDrive' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'date' } },
          type: {
            kind: 'NonNullType',
            type: {
              kind: 'NamedType',
              name: { kind: 'Name', value: 'DateTime' },
            },
          },
        },
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'distance' },
          },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'Int' } },
          },
        },
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'memo' } },
          type: { kind: 'NamedType', name: { kind: 'Name', value: 'String' } },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'recordDrivingRecord' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'date' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'date' },
                },
              },
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'odometerValue' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'distance' },
                },
              },
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'memo' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'memo' },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<RecordDriveMutation, RecordDriveMutationVariables>;
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
            name: { kind: 'Name', value: 'lastReport' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'TotalDistance' },
                },
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'LastOdometerValue' },
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
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'ChartDataSeries' },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TotalDistance' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
          { kind: 'Field', name: { kind: 'Name', value: 'reportDate' } },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'LastOdometerValue' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'ChartDataSeries' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DrivingRecordsConnection' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'nodes' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'odometerValue' },
                },
                { kind: 'Field', name: { kind: 'Name', value: 'recordedAt' } },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'tripDistance' },
                },
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
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'TotalDistance' },
                },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'dailyReports' },
                  selectionSet: {
                    kind: 'SelectionSet',
                    selections: [
                      {
                        kind: 'FragmentSpread',
                        name: { kind: 'Name', value: 'ChartDataSeries' },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TotalDistance' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DistanceReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'odometerValue' } },
          { kind: 'Field', name: { kind: 'Name', value: 'reportDate' } },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'ChartDataSeries' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DrivingRecordsConnection' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'nodes' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'odometerValue' },
                },
                { kind: 'Field', name: { kind: 'Name', value: 'recordedAt' } },
                {
                  kind: 'Field',
                  name: { kind: 'Name', value: 'tripDistance' },
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
          { kind: 'Field', name: { kind: 'Name', value: 'year' } },
          { kind: 'Field', name: { kind: 'Name', value: 'month' } },
          {
            kind: 'FragmentSpread',
            name: { kind: 'Name', value: 'TotalDistance' },
          },
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'dailyReports' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'ChartDataSeries' },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<MonthReportQuery, MonthReportQueryVariables>;
