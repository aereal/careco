import { type CodegenConfig } from '@graphql-codegen/cli';
import { type ClientPresetConfig } from '@graphql-codegen/client-preset';
import { type TypeScriptPluginConfig } from '@graphql-codegen/typescript';

const config: CodegenConfig = {
  ignoreNoDocuments: true,
  schema: './schema.gql',
  documents: ['./frontend/src/**/!(*.gen).{ts,tsx}', './frontend/src/**/*.vue'],
  hooks: {
    afterAllFileWrite: ['pnpm run format'],
  },
  config: {
    scalars: {
      DateTime: {
        input: 'string',
        output: 'Date',
      },
      Month: 'string',
    },
    strictScalars: true,
    defaultScalarType: 'unknown',
    enumsAsTypes: true,
    useTypeImports: true,
    immutableTypes: true,
    skipTypename: true,
  } satisfies TypeScriptPluginConfig,
  generates: {
    'frontend/src/graphql/': {
      preset: 'client',
      presetConfig: {
        fragmentMasking: { unmaskFunctionName: 'getFragmentData' },
      } satisfies ClientPresetConfig,
    },
  },
};

export default config;
