import prettier from 'eslint-config-prettier';
import path from 'node:path';
import js from '@eslint/js';
import paraglide from '@inlang/paraglide-js/eslint';
import svelte from 'eslint-plugin-svelte';
import { defineConfig, includeIgnoreFile } from 'eslint/config';
import globals from 'globals';
import ts from 'typescript-eslint';
import svelteConfig from './svelte.config.js';

// The repo's ignore rules live in the root .gitignore (web/ has none).
const gitignorePath = path.resolve(import.meta.dirname, '../.gitignore');

export default defineConfig(
	includeIgnoreFile(gitignorePath),
	// generated: schema.ts by `make bindings`, src/lib/paraglide/ by the Paraglide Vite plugin
	{ ignores: ['src/lib/api/schema.ts', 'src/lib/paraglide/**'] },
	js.configs.recommended,
	ts.configs.recommended,
	svelte.configs.recommended,
	prettier,
	svelte.configs.prettier,
	{
		languageOptions: { globals: { ...globals.browser, ...globals.node } },
		rules: {
			// typescript-eslint strongly recommend that you do not use the no-undef lint rule on TypeScript projects.
			// see: https://typescript-eslint.io/troubleshooting/faqs/eslint/#i-get-errors-from-the-no-undef-rule-about-global-variables-not-being-defined-even-though-there-are-no-typescript-errors
			'no-undef': 'off'
		}
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
		languageOptions: {
			parserOptions: {
				projectService: true,
				extraFileExtensions: ['.svelte'],
				parser: ts.parser,
				svelteConfig
			}
		}
	},
	{
		// Type-aware rules for app code: catches a forgotten `await` on API calls.
		files: ['src/**/*.ts', 'src/**/*.svelte'],
		languageOptions: { parserOptions: { projectService: true } },
		rules: { '@typescript-eslint/no-floating-promises': 'error' }
	},
	{
		// m[key] defeats Paraglide's tree-shaking (and the typed keys); call m.key() directly.
		files: ['src/**/*.ts', 'src/**/*.svelte'],
		plugins: { paraglide },
		rules: {
			'paraglide/no-dynamic-message-access': ['error', { messageModuleDirectories: ['paraglide'] }]
		}
	},
	{
		// Override or add rule settings here, such as:
		// 'svelte/button-has-type': 'error'
		rules: {}
	}
);
