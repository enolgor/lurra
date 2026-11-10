import { render, screen, waitFor, within } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { searchParams } from 'sv-router';
import { afterEach, describe, expect, test } from 'vitest';
import { json, stubApi } from '$lib/test/api';
import { TodoList, type Todo } from '$lib/todos.svelte';
import Home from './Home.svelte';

const milk: Todo = { id: 1, title: 'buy milk', done: false };
const bread: Todo = { id: 2, title: 'bake bread', done: true };
const serverError = { code: 500, message: 'x', request_id: 'host/abc-000001' };

/** Renders the page with an API that answers with the given responses, in order. */
function setup(...responses: (Response | Promise<Response>)[]) {
	const { client, fetch, request } = stubApi(...responses);
	const todos = new TodoList(client);
	render(Home, { todos });
	return { user: userEvent.setup(), fetch, todos, request };
}
/** A response that arrives when `answer` is called. */
function later() {
	let answer!: (response: Response) => void;
	const response = new Promise<Response>((resolve) => (answer = resolve));
	return { response, answer };
}
const item = (title: string) => screen.getByText(title).closest('li')!;
const loaded = () => screen.findByRole('list');

afterEach(() => {
	searchParams.delete('filter', undefined, { replace: true });
});

describe('loading', () => {
	test('shows a spinner, then the list', async () => {
		const { response, answer } = later();
		setup(response);

		expect(screen.getByRole('status')).toHaveTextContent('Loading…');
		answer(json([milk, bread]));

		await loaded();
		expect(screen.queryByRole('status')).not.toBeInTheDocument();
		expect(screen.getAllByRole('checkbox').map((box) => box.getAttribute('aria-label'))).toEqual([
			'buy milk',
			'bake bread'
		]);
		expect(within(item('bake bread')).getByRole('checkbox')).toBeChecked();
	});

	test('shows "Nothing to do yet" for an empty list', async () => {
		setup(json([]));
		expect(await screen.findByText('Nothing to do yet.')).toBeInTheDocument();
		expect(screen.queryByRole('list')).not.toBeInTheDocument();
	});

	test('a failed load shows the error in place of the list, and "Try again" reloads', async () => {
		const { user } = setup(json(serverError, 500), json([milk]));

		const alert = await screen.findByRole('alert');
		expect(alert).toHaveTextContent('Something went wrong on the server.');
		expect(alert).toHaveTextContent('Request ID: host/abc-000001');

		await user.click(within(alert).getByRole('button', { name: 'Try again' }));
		await loaded();
		expect(screen.queryByRole('alert')).not.toBeInTheDocument();
	});
});

describe('adding', () => {
	test('Add is disabled while the title is blank', async () => {
		const { user } = setup(json([]));
		const add = screen.getByRole('button', { name: 'Add' });

		expect(add).toBeDisabled();
		await user.type(screen.getByRole('textbox', { name: 'New todo' }), '   ');
		expect(add).toBeDisabled();
	});

	test('adds the trimmed title and clears the field', async () => {
		const { user, request } = setup(json([]), json({ id: 3, title: 'call mom', done: false }, 201));
		const input = screen.getByRole('textbox', { name: 'New todo' });

		await user.type(input, '  call mom {Enter}');

		await waitFor(() => expect(item('call mom')).toBeInTheDocument());
		expect(input).toHaveValue('');
		expect(await request(1).json()).toEqual({ title: 'call mom' });
	});

	test('keeps the title when adding fails, and shows the error in a toast', async () => {
		const { user } = setup(json([]), json(serverError, 500));
		const input = screen.getByRole('textbox', { name: 'New todo' });

		await user.type(input, 'call mom{Enter}');

		expect(await screen.findByRole('alert')).toHaveTextContent('Something went wrong');
		expect(input).toHaveValue('call mom');
	});
});

describe('changing', () => {
	test('the checkbox toggles through the API', async () => {
		const { user, request } = setup(json([milk]), json({ ...milk, done: true }));
		await loaded();
		const box = within(item('buy milk')).getByRole('checkbox');

		await user.click(box);

		await waitFor(() => expect(box).toBeChecked());
		expect(await request(1).json()).toEqual({ done: true });
	});

	test('a failed toggle leaves the box as saved, shows a toast, and closing it clears it', async () => {
		const { user } = setup(json([milk]), json(serverError, 500));
		await loaded();
		const box = within(item('buy milk')).getByRole('checkbox');

		await user.click(box);

		const toast = await screen.findByRole('alert');
		expect(box).not.toBeChecked();
		await user.click(within(toast).getByRole('button', { name: 'Close' }));
		expect(screen.queryByRole('alert')).not.toBeInTheDocument();
	});

	test("a todo's controls are disabled while its request is in flight", async () => {
		const { response, answer } = later();
		const { user } = setup(json([milk, bread]), response);
		await loaded();

		await user.click(screen.getByRole('button', { name: 'Delete “buy milk”' }));

		const busy = within(item('buy milk'));
		expect(busy.getByRole('checkbox')).toBeDisabled();
		expect(busy.getByRole('button', { name: 'Edit “buy milk”' })).toBeDisabled();
		expect(busy.getByRole('button', { name: 'Delete “buy milk”' })).toBeDisabled();
		expect(within(item('bake bread')).getByRole('checkbox')).toBeEnabled();

		answer(new Response(null, { status: 204 }));
		await waitFor(() => expect(screen.queryByText('buy milk')).not.toBeInTheDocument());
	});
});

describe('editing', () => {
	async function startEditing(...responses: Response[]) {
		const context = setup(json([milk]), ...responses);
		await loaded();
		await context.user.click(screen.getByRole('button', { name: 'Edit “buy milk”' }));
		return { ...context, input: screen.getByRole('textbox', { name: 'Title' }) };
	}

	test('opens a focused field with the title; Enter saves', async () => {
		const { user, input, request } = await startEditing(json({ ...milk, title: 'buy oat milk' }));
		expect(input).toHaveFocus();
		expect(input).toHaveValue('buy milk');

		await user.clear(input);
		await user.type(input, 'buy oat milk{Enter}');

		expect(await screen.findByText('buy oat milk')).toBeInTheDocument();
		expect(screen.queryByRole('textbox', { name: 'Title' })).not.toBeInTheDocument();
		expect(await request(1).json()).toEqual({ title: 'buy oat milk' });
	});

	test('leaving the field saves too', async () => {
		const { user, input } = await startEditing(json({ ...milk, title: 'buy oat milk' }));

		await user.type(input, '{End} now');
		await user.click(document.body);

		expect(await screen.findByText('buy oat milk')).toBeInTheDocument();
	});

	test.each([
		['Escape', (input: HTMLElement) => userEvent.type(input, 'x{Escape}')],
		[
			'a blank title',
			(input: HTMLElement) => userEvent.type(input, '{Control>}a{/Control}{Delete}{Enter}')
		],
		['an unchanged title', (input: HTMLElement) => userEvent.type(input, '{Enter}')]
	])('%s closes the field without saving', async (_, finish) => {
		const { input, fetch } = await startEditing();

		await finish(input);

		expect(screen.queryByRole('textbox', { name: 'Title' })).not.toBeInTheDocument();
		expect(screen.getByText('buy milk')).toBeInTheDocument();
		expect(fetch).toHaveBeenCalledTimes(1); // only the initial load
	});
});

describe('filter', () => {
	const tab = (name: string) => screen.getByRole('tab', { name });

	test('the tabs filter the list and write ?filter=, with "All" as no parameter', async () => {
		const { user } = setup(json([milk, bread]));
		await loaded();
		expect(tab('All')).toHaveAttribute('aria-selected', 'true');

		await user.click(tab('Done'));
		expect(location.search).toBe('?filter=done');
		expect(tab('Done')).toHaveAttribute('aria-selected', 'true');
		expect(screen.queryByText('buy milk')).not.toBeInTheDocument();
		expect(screen.getByText('bake bread')).toBeInTheDocument();

		await user.click(tab('To do'));
		expect(location.search).toBe('?filter=undone');
		expect(screen.getByText('buy milk')).toBeInTheDocument();

		await user.click(tab('All'));
		expect(location.search).toBe('');
		expect(screen.getAllByRole('checkbox')).toHaveLength(2);
	});

	test('reads the filter from the URL, as after a reload', async () => {
		searchParams.set('filter', 'done', { replace: true });
		setup(json([milk, bread]));
		await loaded();

		expect(tab('Done')).toHaveAttribute('aria-selected', 'true');
		expect(screen.queryByText('buy milk')).not.toBeInTheDocument();
	});

	test('an unknown filter in the URL shows all', async () => {
		searchParams.set('filter', 'nope', { replace: true });
		setup(json([milk]));
		await loaded();

		expect(tab('All')).toHaveAttribute('aria-selected', 'true');
	});

	test('says so when the filter matches nothing', async () => {
		searchParams.set('filter', 'done', { replace: true });
		setup(json([milk]));

		expect(await screen.findByText('No todos here.')).toBeInTheDocument();
	});
});
