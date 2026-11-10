import { render, screen } from '@testing-library/svelte';
import { createRawSnippet } from 'svelte';
import { afterEach, expect, test } from 'vitest';
import '$lib/i18n.svelte';
import { setLocale } from '$lib/paraglide/runtime.js';
import type { ApiError } from '$lib/todos.svelte';
import ApiErrorAlert from './ApiErrorAlert.svelte';

afterEach(() => {
	void setLocale('en');
	localStorage.clear();
});

// what the server sends as `message`: English and technical, never shown to users
const serverText = 'server message';

test.each<[string, Omit<ApiError, 'message'>, string]>([
	[
		'the error code over the status',
		{ code: 400, error: 'title_invalid' },
		"The title can't be empty."
	],
	['no connection (code 0)', { code: 0 }, "Can't reach the server. Check your connection."],
	['a 400 without a code', { code: 400 }, "That request isn't valid."],
	['404', { code: 404 }, 'This todo no longer exists.'],
	['413', { code: 413 }, "That's too long."],
	['500', { code: 500 }, 'Something went wrong on the server.'],
	['any other status', { code: 418 }, 'Something went wrong.']
])('translates %s', (_, error, text) => {
	render(ApiErrorAlert, { error: { ...error, message: serverText } });

	const alert = screen.getByRole('alert');
	expect(alert).toHaveTextContent(text);
	expect(alert).not.toHaveTextContent(serverText);
});

test('shows the request ID when there is one', () => {
	const { unmount } = render(ApiErrorAlert, {
		error: { code: 500, message: 'x', request_id: 'host/abc-000001' }
	});
	expect(screen.getByText('Request ID: host/abc-000001')).toBeInTheDocument();
	unmount();

	render(ApiErrorAlert, { error: { code: 500, message: 'x' } });
	expect(screen.queryByText(/Request ID/)).not.toBeInTheDocument();
});

test('renders an action next to the message', () => {
	render(ApiErrorAlert, {
		error: { code: 500, message: 'x' },
		action: createRawSnippet(() => ({ render: () => '<button type="button">Try again</button>' }))
	});

	expect(screen.getByRole('button', { name: 'Try again' })).toBeInTheDocument();
});

test('follows the language live', async () => {
	render(ApiErrorAlert, { error: { code: 404, message: 'x', request_id: 'r1' } });

	void setLocale('de');

	expect(await screen.findByText('Diese Aufgabe gibt es nicht mehr.')).toBeInTheDocument();
	expect(screen.getByText('Anfrage-ID: r1')).toBeInTheDocument();
});
