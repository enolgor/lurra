// Adds DOM matchers such as toBeInTheDocument() to expect().
import '@testing-library/jest-dom/vitest';

// jsdom has no IntersectionObserver; sv-router uses it to preload visible links.
globalThis.IntersectionObserver ??= class {
	observe() {}
	unobserve() {}
	disconnect() {}
	takeRecords() {
		return [];
	}
} as unknown as typeof IntersectionObserver;
