import { describe, expect, test } from 'vitest';
import { localeFlag, localeName } from './locale';

describe('localeName', () => {
	test.each([
		['en', 'English'],
		['es', 'Español'],
		['de', 'Deutsch'],
		['ja', '日本語']
	])('%s → %s', (locale, name) => {
		expect(localeName(locale)).toBe(name);
	});
});

describe('localeFlag', () => {
	test.each([
		['en', '🇺🇸'],
		['es', '🇪🇸'],
		['de', '🇩🇪'],
		['pt-BR', '🇧🇷']
	])('%s → %s', (locale, flag) => {
		expect(localeFlag(locale)).toBe(flag);
	});

	test('no flag when the locale has no country', () => {
		expect(localeFlag('eo')).toBe(''); // Esperanto: region 001 (world)
	});
});
