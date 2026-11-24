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
  tireID: string | number;
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

export type TireOptionFragment = {
  readonly id: string;
  readonly name: string;
} & { ' $fragmentName'?: 'TireOptionFragment' };

export type LastTireFragment = {
  readonly tire: { readonly id: string } | null;
} & { ' $fragmentName'?: 'LastTireFragment' };

export type ListTiresQueryVariables = Exact<{
  includeRetired: boolean;
}>;

export type ListTiresQuery = {
  readonly tires: ReadonlyArray<{
    readonly id: string;
    readonly name: string;
    readonly purchasedOn: Date | null;
    readonly retired: boolean;
  }>;
};

export type AddTireMutationVariables = Exact<{
  name: string;
  purchasedOn?: string | null | undefined;
}>;

export type AddTireMutation = { readonly addTire: { readonly id: string } };

export type RetireTireMutationVariables = Exact<{
  id: string | number;
}>;

export type RetireTireMutation = {
  readonly retireTire: { readonly id: string };
};

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
      LastTireFragment: LastTireFragment;
    };
  } | null;
  readonly tires: ReadonlyArray<{
    ' $fragmentRefs'?: { TireOptionFragment: TireOptionFragment };
  }>;
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
export const TireOptionFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TireOption' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'Tire' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'id' } },
          { kind: 'Field', name: { kind: 'Name', value: 'name' } },
        ],
      },
    },
  ],
} as unknown as DocumentNode<TireOptionFragment, unknown>;
export const LastTireFragmentDoc = {
  kind: 'Document',
  definitions: [
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'LastTire' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DailyReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'tire' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'id' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<LastTireFragment, unknown>;
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
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'tireID' },
          },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'ID' } },
          },
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
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'tireID' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'tireID' },
                },
              },
            ],
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<RecordDriveMutation, RecordDriveMutationVariables>;
export const ListTiresDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'query',
      name: { kind: 'Name', value: 'ListTires' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'includeRetired' },
          },
          type: {
            kind: 'NonNullType',
            type: {
              kind: 'NamedType',
              name: { kind: 'Name', value: 'Boolean' },
            },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'tires' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'includeRetired' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'includeRetired' },
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'id' } },
                { kind: 'Field', name: { kind: 'Name', value: 'name' } },
                { kind: 'Field', name: { kind: 'Name', value: 'purchasedOn' } },
                { kind: 'Field', name: { kind: 'Name', value: 'retired' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<ListTiresQuery, ListTiresQueryVariables>;
export const AddTireDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'AddTire' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'name' } },
          type: {
            kind: 'NonNullType',
            type: {
              kind: 'NamedType',
              name: { kind: 'Name', value: 'String' },
            },
          },
        },
        {
          kind: 'VariableDefinition',
          variable: {
            kind: 'Variable',
            name: { kind: 'Name', value: 'purchasedOn' },
          },
          type: {
            kind: 'NamedType',
            name: { kind: 'Name', value: 'DateTime' },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'addTire' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'name' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'name' },
                },
              },
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'purchasedOn' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'purchasedOn' },
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'id' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<AddTireMutation, AddTireMutationVariables>;
export const RetireTireDocument = {
  kind: 'Document',
  definitions: [
    {
      kind: 'OperationDefinition',
      operation: 'mutation',
      name: { kind: 'Name', value: 'RetireTire' },
      variableDefinitions: [
        {
          kind: 'VariableDefinition',
          variable: { kind: 'Variable', name: { kind: 'Name', value: 'id' } },
          type: {
            kind: 'NonNullType',
            type: { kind: 'NamedType', name: { kind: 'Name', value: 'ID' } },
          },
        },
      ],
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'retireTire' },
            arguments: [
              {
                kind: 'Argument',
                name: { kind: 'Name', value: 'id' },
                value: {
                  kind: 'Variable',
                  name: { kind: 'Name', value: 'id' },
                },
              },
            ],
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'id' } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<RetireTireMutation, RetireTireMutationVariables>;
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
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'LastTire' },
                },
              ],
            },
          },
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'tires' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                {
                  kind: 'FragmentSpread',
                  name: { kind: 'Name', value: 'TireOption' },
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
      name: { kind: 'Name', value: 'LastTire' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'DailyReport' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          {
            kind: 'Field',
            name: { kind: 'Name', value: 'tire' },
            selectionSet: {
              kind: 'SelectionSet',
              selections: [
                { kind: 'Field', name: { kind: 'Name', value: 'id' } },
              ],
            },
          },
        ],
      },
    },
    {
      kind: 'FragmentDefinition',
      name: { kind: 'Name', value: 'TireOption' },
      typeCondition: {
        kind: 'NamedType',
        name: { kind: 'Name', value: 'Tire' },
      },
      selectionSet: {
        kind: 'SelectionSet',
        selections: [
          { kind: 'Field', name: { kind: 'Name', value: 'id' } },
          { kind: 'Field', name: { kind: 'Name', value: 'name' } },
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
