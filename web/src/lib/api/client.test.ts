import { expect, test, vi } from 'vitest';
import { api } from './client';

function respond(status: number, body: unknown) {
	return vi.fn<typeof fetch>(async () => Response.json(body, { status }));
}

test('sends requests under /api/v1', async () => {
	const fetch = respond(200, [{ id: 1, title: 'buy milk', done: false }]);

	const { data, error } = await api.GET('/todos', { fetch });

	expect(error).toBeUndefined();
	expect(data).toEqual([{ id: 1, title: 'buy milk', done: false }]);
	const req = fetch.mock.calls[0][0] as Request;
	expect(req.method).toBe('GET');
	expect(new URL(req.url).pathname).toBe('/api/v1/todos');
});

test('returns API errors with their request_id', async () => {
	const body = { code: 404, message: 'todo not found', request_id: 'host/abc-000001' };
	const fetch = respond(404, body);

	const { data, error, response } = await api.GET('/todos/{id}', {
		params: { path: { id: 999 } },
		fetch
	});

	expect(data).toBeUndefined();
	expect(response.status).toBe(404);
	expect(error).toEqual(body);
});
