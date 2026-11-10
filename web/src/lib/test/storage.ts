import { vi } from 'vitest';

/** Makes localStorage's `method` throw, as when the browser blocks storage (restore with vi.restoreAllMocks). */
export const blockStorage = (method: 'getItem' | 'setItem') =>
	vi.spyOn(Storage.prototype, method).mockImplementation(() => {
		throw new DOMException('blocked', 'SecurityError');
	});
