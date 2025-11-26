import { type NextConfig } from 'next';

const config: NextConfig = {
  typedRoutes: true,
  experimental: {
    // PlaywrightでServer Componentのfetchをモックするために必要
    testProxy: true,
  },
};

export default config;
