import type { ApiError, Todo } from '$lib/todos.svelte';

/** Two todos: one to do, one done. */
export const milk: Todo = { id: 1, title: 'buy milk', done: false };
export const bread: Todo = { id: 2, title: 'bake bread', done: true };

/** A 500 as the API sends it, with the request ID shown for bug reports. */
export const serverError: ApiError = { code: 500, message: 'x', request_id: 'host/abc-000001' };
