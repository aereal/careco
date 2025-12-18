export const getEventNumberValue = <E extends Event>(
  e: E,
): number | undefined => {
  if (!(e.target instanceof HTMLInputElement)) {
    return undefined;
  }
  const num = e.target.valueAsNumber;
  if (isNaN(num)) {
    return undefined;
  }
  return num;
};
