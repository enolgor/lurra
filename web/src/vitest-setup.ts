// Adds DOM matchers such as toBeInTheDocument() to expect().
import '@testing-library/jest-dom/vitest';
import { afterEach, beforeEach, vi, type MockInstance } from 'vitest';

// A test fails if anything logs a warning or an error (Svelte's runtime warnings, jsdom's "Not
// implemented", a library's console.error). The output is still printed, to see what it was.
let spies: MockInstance[] = [];
beforeEach(() => {
	spies = [vi.spyOn(console, 'warn'), vi.spyOn(console, 'error')];
});
afterEach(() => {
	const logged = spies.flatMap((spy) => spy.mock.calls.map((args) => args.join(' ')));
	for (const spy of spies) spy.mockRestore();
	if (logged.length > 0)
		throw new Error(`the test logged a warning or error:\n${logged.join('\n')}`);
});

// jsdom has no IntersectionObserver; sv-router uses it to preload visible links.
globalThis.IntersectionObserver ??= class {
	observe() {}
	unobserve() {}
	disconnect() {}
	takeRecords() {
		return [];
	}
} as unknown as typeof IntersectionObserver;

// jsdom's scrollTo only logs "Not implemented"; sv-router scrolls to the top on navigation.
window.scrollTo = () => {};
