export const register = async (): Promise<void> => {
  const runtime = process.env.NEXT_RUNTIME;
  switch (runtime) {
    case 'nodejs':
      (await import('./instrumentation.node')).register();
      return;
    default:
      console.log(`-> unknown NEXT_RUNTIME: ${runtime}`);
  }
};
