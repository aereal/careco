/* eslint-disable */
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
  Month: { input: string; output: string };
};

export type DailyReport = {
  readonly day: Scalars['Int']['output'];
  readonly distanceKillos: Scalars['Int']['output'];
  readonly month: Scalars['Month']['output'];
  readonly year: Scalars['Int']['output'];
};

export type DrivingRecord = {
  readonly distanceKillos: Scalars['Int']['output'];
  readonly memo?: Maybe<Scalars['String']['output']>;
  readonly recordedAt: Scalars['DateTime']['output'];
};

export type MonthlyReport = {
  readonly dailyStatistics: ReadonlyArray<DailyReport>;
  readonly distanceKillos: Scalars['Int']['output'];
  readonly month: Scalars['Month']['output'];
  readonly year: Scalars['Int']['output'];
};

export type Mutation = {
  readonly recordDrivingRecord: Scalars['Boolean']['output'];
};

export type MutationRecordDrivingRecordArgs = {
  date: Scalars['DateTime']['input'];
  distanceKillos: Scalars['Int']['input'];
  memo?: InputMaybe<Scalars['String']['input']>;
};

export type Query = {
  readonly monthlyReport: MonthlyReport;
  readonly recentDrivingRecords: ReadonlyArray<DrivingRecord>;
  readonly totalStatistics: TotalStatistics;
  readonly yearlyReport: YearlyReport;
};

export type QueryMonthlyReportArgs = {
  month: Scalars['Month']['input'];
  year: Scalars['Int']['input'];
};

export type QueryRecentDrivingRecordsArgs = {
  first: Scalars['Int']['input'];
};

export type QueryYearlyReportArgs = {
  year: Scalars['Int']['input'];
};

export type TotalStatistics = {
  readonly distanceKillos: Scalars['Int']['output'];
};

export type YearlyReport = {
  readonly distanceKillos: Scalars['Int']['output'];
  readonly monthlyStatistics: ReadonlyArray<MonthlyReport>;
  readonly year: Scalars['Int']['output'];
};
