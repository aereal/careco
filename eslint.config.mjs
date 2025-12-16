import vitest from '@vitest/eslint-plugin';
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting';
import {
  defineConfigWithVueTs,
  vueTsConfigs,
} from '@vue/eslint-config-typescript';
import playwright from 'eslint-plugin-playwright';
import vue from 'eslint-plugin-vue';
import { globalIgnores } from 'eslint/config';

/**
 * @type {Array<import('eslint').Linter.Config<import('eslint').Linter.RulesRecord>>}
 */
const eslintConfig = defineConfigWithVueTs(
  {
    files: ['frontend/src/**/*.{ts,mts,vue}'],
  },
  globalIgnores([
    'node_modules/**',
    'dist/**/*',
    'frontend/src/graphql/**/*.ts',
  ]),
  vue.configs['flat/essential'],
  vueTsConfigs.recommended,
  {
    ...vitest.configs.recommended,
    files: ['frontend/src/**/*.{test,spec}.ts'],
  },
  { ...playwright.configs['flat/recommended'], files: ['frontend/e2e/**/*'] },
  {
    rules: {
      'vue/multi-word-component-names': ['off'],
    },
  },
  skipFormatting,
);

export default eslintConfig;
