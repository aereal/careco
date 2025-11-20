export const date = '日付' as const;
export const distance = '走行距離' as const;
export const totalDistance = '総走行距離' as const;

export type KeyDate = typeof date;
export type KeyDistance = typeof distance;
export type KeyTotalDistance = typeof totalDistance;
