/** The value, or a failed test when it's null or undefined (e.g. an element that isn't there). */
export function defined<T>(value: T | null | undefined): T {
	if (value == null) throw new Error(`expected a value, got ${String(value)}`);
	return value;
}
