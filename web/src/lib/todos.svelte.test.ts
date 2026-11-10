import { describe, expect, test } from 'vitest';
import { json, stubApi } from '$lib/test/api';
import { bread, milk, serverError } from '$lib/test/todos';
import { TodoList } from '$lib/todos.svelte';

const titleInvalid = {
	code: 400,
	message: 'title is empty',
	error: 'title_invalid' as const,
	request_id: 'host/abc-000001'
};

/** A TodoList whose API answers with the given responses, in order. */
function setup(...responses: Response[]) {
	const { client, fetch, request } = stubApi(...responses);
	return { todos: new TodoList(client), fetch, request };
}
const noContent = () => new Response(null, { status: 204 });

describe('load', () => {
	test('loads the list, loading while it waits', async () => {
		const { todos, request } = setup(json([milk, bread]));

		const loading = todos.load();
		expect(todos.loading).toBe(true);
		await loading;

		expect(todos.loading).toBe(false);
		expect(todos.items).toEqual([milk, bread]);
		expect(todos.loadError).toBeNull();
		expect(request(0).method).toBe('GET');
		expect(new URL(request(0).url).pathname).toBe('/api/v1/todos');
	});

	test('keeps the API error for the inline alert', async () => {
		const { todos } = setup(json(serverError, 500));

		await todos.load();

		expect(todos.loading).toBe(false);
		expect(todos.loadError).toEqual(serverError);
		expect(todos.actionError).toBeNull();
		expect(todos.items).toEqual([]);
	});

	test('records a network failure as code 0', async () => {
		const { todos, fetch } = setup();
		fetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

		await todos.load();

		expect(todos.loadError).toEqual({ code: 0, message: 'TypeError: Failed to fetch' });
	});

	test('a new load clears the previous error', async () => {
		const { todos } = setup(json({ code: 500, message: 'x' }, 500), json([milk]));

		await todos.load();
		await todos.load();

		expect(todos.loadError).toBeNull();
		expect(todos.items).toEqual([milk]);
	});
});

describe('changes', () => {
	test('add appends the created todo', async () => {
		const created = { id: 3, title: 'call mom', done: false };
		const { todos, request } = setup(json([milk]), json(created, 201));
		await todos.load();

		expect(await todos.add('call mom')).toBe(true);

		expect(todos.items).toEqual([milk, created]);
		expect(request(1).method).toBe('POST');
		expect(await request(1).json()).toEqual({ title: 'call mom' });
	});

	test('toggle replaces the todo with the updated one', async () => {
		const { todos, request } = setup(json([milk, bread]), json({ ...milk, done: true }));
		await todos.load();

		await todos.toggle(milk);

		expect(todos.items).toEqual([{ ...milk, done: true }, bread]);
		expect(request(1).method).toBe('PATCH');
		expect(new URL(request(1).url).pathname).toBe('/api/v1/todos/1');
		expect(await request(1).json()).toEqual({ done: true });
	});

	test('rename replaces the todo with the updated one', async () => {
		const { todos, request } = setup(json([milk]), json({ ...milk, title: 'buy oat milk' }));
		await todos.load();

		await todos.rename(milk, 'buy oat milk');

		expect(todos.items).toEqual([{ ...milk, title: 'buy oat milk' }]);
		expect(await request(1).json()).toEqual({ title: 'buy oat milk' });
	});

	test('remove drops the todo', async () => {
		const { todos, request } = setup(json([milk, bread]), noContent());
		await todos.load();

		await todos.remove(milk);

		expect(todos.items).toEqual([bread]);
		expect(request(1).method).toBe('DELETE');
		expect(new URL(request(1).url).pathname).toBe('/api/v1/todos/1');
	});

	test.each([
		['add', (todos: TodoList) => todos.add('')],
		['toggle', (todos: TodoList) => todos.toggle(milk)],
		['rename', (todos: TodoList) => todos.rename(milk, '')],
		['remove', (todos: TodoList) => todos.remove(milk)]
	])('a failed %s changes nothing and keeps the error for the toast', async (_, change) => {
		const { todos } = setup(json([milk]), json(titleInvalid, 400));
		await todos.load();

		expect(await change(todos)).toBe(false);

		expect(todos.items).toEqual([milk]);
		expect(todos.actionError).toEqual(titleInvalid);
		expect(todos.loadError).toBeNull();
	});

	test('dismissError clears the toast error', async () => {
		const { todos } = setup(json(titleInvalid, 400));
		await todos.add('');

		todos.dismissError();

		expect(todos.actionError).toBeNull();
	});
});

describe('busy', () => {
	test.each([
		['toggle', (todos: TodoList) => todos.toggle(milk)],
		['rename', (todos: TodoList) => todos.rename(milk, 'x')],
		['remove', (todos: TodoList) => todos.remove(milk)]
	])('%s marks the todo busy until the response, success or not', async (_, change) => {
		const { todos, fetch } = setup(json([milk]));
		await todos.load();
		let answer!: (response: Response) => void;
		fetch.mockReturnValueOnce(new Promise((resolve) => (answer = resolve)));

		const pending = change(todos);
		expect(todos.busy.has(milk.id)).toBe(true);
		answer(json(titleInvalid, 400));
		await pending;

		expect(todos.busy.has(milk.id)).toBe(false);
	});
});

describe('filter', () => {
	test.each([
		['all', [milk, bread]],
		['done', [bread]],
		['undone', [milk]]
	] as const)('%s shows the matching todos', async (filter, visible) => {
		const { todos } = setup(json([milk, bread]));
		await todos.load();

		todos.filter = filter;

		expect(todos.visible).toEqual(visible);
	});
});
