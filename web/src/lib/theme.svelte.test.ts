// node:fs reads public/theme.js; Node types only here, app code stays browser-only
/// <reference types="node" />
import { readFileSync } from 'node:fs';
import { afterEach, describe, expect, test, vi } from 'vitest';
import { blockStorage } from '$lib/test/storage';
import { detectTheme, getTheme, setTheme, themeStorageKey } from '$lib/theme.svelte';

const html = document.documentElement;

afterEach(() => {
	vi.restoreAllMocks();
	setTheme('system');
});

describe('detectTheme', () => {
	test('uses the saved choice', () => {
		localStorage.setItem(themeStorageKey, 'dark');
		expect(detectTheme()).toBe('dark');
	});

	test('follows the system without a valid saved choice', () => {
		expect(detectTheme()).toBe('system');
		localStorage.setItem(themeStorageKey, 'neon');
		expect(detectTheme()).toBe('system');
	});

	test('follows the system when storage is blocked', () => {
		blockStorage('getItem');
		expect(detectTheme()).toBe('system');
	});
});

describe('setTheme', () => {
	test('applies and saves light or dark', () => {
		setTheme('dark');

		expect(getTheme()).toBe('dark');
		expect(html.dataset.theme).toBe('dark');
		expect(localStorage.getItem(themeStorageKey)).toBe('dark');
	});

	test('system removes data-theme and forgets the saved choice', () => {
		setTheme('light');
		setTheme('system');

		expect(getTheme()).toBe('system');
		expect(html).not.toHaveAttribute('data-theme');
		expect(localStorage.getItem(themeStorageKey)).toBeNull();
	});

	test('still applies when storage is blocked', () => {
		blockStorage('setItem');
		setTheme('dark');

		expect(getTheme()).toBe('dark');
		expect(html.dataset.theme).toBe('dark');
	});
});

// public/theme.js runs before the app; it must read the same key and apply only light or dark.
describe('public/theme.js', () => {
	const script = readFileSync('public/theme.js', 'utf8'); // tests run from web/
	// run it as the browser does: a <script> element (jsdom executes it)
	const run = () => {
		const element = document.createElement('script');
		element.textContent = script;
		document.head.append(element);
		element.remove();
	};

	test.each(['light', 'dark'])('applies a saved %s before the app loads', (theme) => {
		localStorage.setItem(themeStorageKey, theme);
		run();
		expect(html.dataset.theme).toBe(theme);
	});

	test('leaves the system theme without a valid saved choice', () => {
		run();
		expect(html).not.toHaveAttribute('data-theme');
		localStorage.setItem(themeStorageKey, 'system');
		run();
		expect(html).not.toHaveAttribute('data-theme');
	});

	test('does not fail when storage is blocked', () => {
		const errors: ErrorEvent[] = [];
		const onError = (event: ErrorEvent) => errors.push(event);
		window.addEventListener('error', onError);
		blockStorage('getItem');

		run();

		window.removeEventListener('error', onError);
		expect(errors).toEqual([]);
		expect(html).not.toHaveAttribute('data-theme');
	});
});
