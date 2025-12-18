/**
 * @type {import('prettier').Config}
 */
const config = {
  arrowParens: 'always',
  bracketSpacing: true,
  jsxSingleQuote: true,
  semi: true,
  singleQuote: true,
  trailingComma: 'all',
  overrides: [
    { files: ['frontend/**/*'], excludeFiles: ['typed-router.d.ts'] },
  ],
};

export default config;
