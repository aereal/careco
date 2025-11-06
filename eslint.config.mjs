import coreWebVitals from 'eslint-config-next/core-web-vitals';
import ts from 'eslint-config-next/typescript';
import { defineConfig, globalIgnores } from 'eslint/config';

/**
 * @type {Array<import('eslint').Linter.Config<import('eslint').Linter.RulesRecord>>}
 */
const eslintConfig = defineConfig([
  ...coreWebVitals,
  ...ts,
  globalIgnores([
    'node_modules/**',
    '.next/**',
    'out/**',
    'next-env.d.ts',
    'src/graphql/**/*.ts',
  ]),
]);

export default eslintConfig;
