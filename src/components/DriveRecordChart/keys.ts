export const date = '日付' as const;
export const tripDistance = '走行距離' as const;
export const odometerValue = '総走行距離' as const;

export type KeyDate = typeof date;
export type KeyTripDistance = typeof tripDistance;
export type KeyOdometerValue = typeof odometerValue;
