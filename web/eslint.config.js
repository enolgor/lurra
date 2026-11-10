import prettier from 'eslint-config-prettier';
import path from 'node:path';
import js from '@eslint/js';
import paraglide from '@inlang/paraglide-js/eslint';
import betterTailwindcss from 'eslint-plugin-better-tailwindcss';
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
	// Type-aware rules (they read the TypeScript types): forgotten `await`s, promises nobody
	// handles, conditions that are always true, unsafe `any`, needless `!`.
	ts.configs.strictTypeChecked,
	svelte.configs.recommended,
	prettier,
	svelte.configs.prettier,
	{
		languageOptions: {
			globals: { ...globals.browser, ...globals.node },
			// the types come from each file's tsconfig (tsconfig.app.json, tsconfig.node.json)
			parserOptions: { projectService: true, tsconfigRootDir: import.meta.dirname }
		},
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
				extraFileExtensions: ['.svelte'],
				parser: ts.parser,
				svelteConfig
			}
		}
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
		// Tailwind class checks as errors: canonical names (the editor's suggestCanonicalClasses),
		// unknown, conflicting, duplicate and deprecated classes. Class order and line wrapping are
		// Prettier's job (prettier-plugin-tailwindcss), so those two rules are off.
		files: ['src/**/*.svelte', 'src/**/*.ts'],
		extends: [betterTailwindcss.configs['recommended-error']],
		settings: { 'better-tailwindcss': { entryPoint: 'src/app.css' } },
		rules: {
			'better-tailwindcss/enforce-consistent-class-order': 'off',
			'better-tailwindcss/enforce-consistent-line-wrapping': 'off'
		}
	},
	{
		// Override or add rule settings here, such as:
		// 'svelte/button-has-type': 'error'
		rules: {
			// onclick={() => todos.remove(todo)}: handlers return nothing anyway
			'@typescript-eslint/no-confusing-void-expression': ['error', { ignoreArrowShorthand: true }],
			// numbers in messages (`request ${n}`) are fine
			'@typescript-eslint/restrict-template-expressions': ['error', { allowNumber: true }]
		}
	},
	{
		// Plain JS outside any tsconfig (the tool configs, public/theme.js): no types to check.
		// Last, so no rule setting above turns a type-aware rule back on for these files.
		files: ['*.js', 'public/*.js'],
		extends: [ts.configs.disableTypeChecked]
	}
);
