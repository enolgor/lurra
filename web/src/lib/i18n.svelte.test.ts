import { afterEach, describe, expect, test, vi } from 'vitest';
import { getLocale, localStorageKey, setLocale } from '$lib/paraglide/runtime.js';
import { m } from '$lib/paraglide/messages.js';
import { detectLocale } from './i18n.svelte';

const blockStorage = (method: 'getItem' | 'setItem') =>
	vi.spyOn(Storage.prototype, method).mockImplementation(() => {
		throw new DOMException('blocked', 'SecurityError');
	});
const browserLanguages = (languages: string[]) =>
	vi.spyOn(navigator, 'languages', 'get').mockReturnValue(languages);

afterEach(() => {
	vi.restoreAllMocks();
	void setLocale('en');
	localStorage.clear();
});

describe('detectLocale', () => {
	test('prefers the saved choice', () => {
		localStorage.setItem(localStorageKey, 'es');
		browserLanguages(['de']);
		expect(detectLocale()).toBe('es');
	});

	test('ignores an unknown saved value and uses the browser language', () => {
		localStorage.setItem(localStorageKey, 'xx');
		browserLanguages(['fr-FR', 'de-AT']);
		expect(detectLocale()).toBe('de');
	});

	test('falls back to the base locale', () => {
		browserLanguages(['fr', 'it']);
		expect(detectLocale()).toBe('en');
	});

	test('uses the browser language when storage is blocked', () => {
		blockStorage('getItem');
		browserLanguages(['es-MX']);
		expect(detectLocale()).toBe('es');
	});
});

describe('setLocale', () => {
	test('switches messages live, updates <html> and saves the choice', () => {
		void setLocale('de');

		expect(getLocale()).toBe('de');
		expect(m.language_selector_trigger()).toBe('Sprache');
		expect(document.documentElement.lang).toBe('de');
		expect(document.documentElement.dir).toBe('ltr');
		expect(localStorage.getItem(localStorageKey)).toBe('de');
	});

	test('still switches when storage is blocked', () => {
		blockStorage('setItem');
		void setLocale('es');

		expect(getLocale()).toBe('es');
		expect(document.documentElement.lang).toBe('es');
	});
});
