export const getEventDateValue = <E extends Event>(e: E): Date | null => {
  if (!(e.target instanceof HTMLInputElement)) {
    return null;
  }
  return e.target.valueAsDate;
};
