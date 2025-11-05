import { type CodegenConfig } from '@graphql-codegen/cli';
import { type ClientPresetConfig } from '@graphql-codegen/client-preset';
import { type TypeScriptPluginConfig } from '@graphql-codegen/typescript';

const config: CodegenConfig = {
  ignoreNoDocuments: true,
  schema: './schema.gql',
  documents: ['./src/**/!(*.gen).{ts,tsx}'],
  hooks: {
    afterAllFileWrite: ['prettier --write'],
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
    'src/graphql/': {
      preset: 'client',
      presetConfig: {
        fragmentMasking: { unmaskFunctionName: 'getFragmentData' },
      } satisfies ClientPresetConfig,
    },
  },
};

export default config;
