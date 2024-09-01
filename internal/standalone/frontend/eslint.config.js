import globals from 'globals'
import pluginJs from '@eslint/js'
import tseslint from 'typescript-eslint'
import solid from 'eslint-plugin-solid/configs/typescript'
import eslintConfigPrettier from 'eslint-config-prettier'
import * as tsParser from '@typescript-eslint/parser'

export default [
    {
        ignores: ['**/wailsjs/**', '**/dist/**'],
    },
    {
        files: ['**/*.{js,mjs,cjs,ts}'],
    },
    {
        languageOptions: {
            globals: globals.browser,
        },
    },
    pluginJs.configs.recommended,
    ...tseslint.configs.recommended,
    {
        files: ['**/*.{ts,tsx}'],
        ...solid,
        languageOptions: {
            parser: tsParser,
            parserOptions: {
                project: 'tsconfig.json',
            },
        },
    },
    eslintConfigPrettier,
]
