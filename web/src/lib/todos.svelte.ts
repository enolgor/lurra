// The todo list: loads it, changes it through the API, and keeps it in $state. The list is only
// updated from API responses, so a failed change is simply not applied.
import type { Client } from 'openapi-fetch';
import { SvelteSet } from 'svelte/reactivity';
import { api } from '$lib/api/client';
import type { components, paths } from '$lib/api/schema';

export type Todo = components['schemas']['model.Todo'];
/** An API error response; `code: 0` means the request didn't reach the server. */
export type ApiError = components['schemas']['httpx.HTTPError'];

export const filters = ['all', 'done', 'undone'] as const;
export type Filter = (typeof filters)[number];

type Result<T> = { data?: T; error?: ApiError };

export class TodoList {
	items = $state<Todo[]>([]);
	loading = $state(false);
	/** Why the list couldn't load (shown in place of the list). */
	loadError = $state<ApiError | null>(null);
	/** Why the last change failed (shown in a toast until dismissed). */
	actionError = $state<ApiError | null>(null);
	filter = $state<Filter>('all');
	/** Ids of todos with a request in flight; their controls are disabled until it answers. */
	busy = new SvelteSet<number>();

	/** The items the current filter shows. */
	visible = $derived(
		this.filter === 'all'
			? this.items
			: this.items.filter((todo) => todo.done === (this.filter === 'done'))
	);

	#client: Client<paths>;

	/** Uses the app's API client; tests pass one with a stubbed fetch. */
	constructor(client: Client<paths> = api) {
		this.#client = client;
	}

	async load() {
		this.loading = true;
		this.loadError = null;
		const result = await this.#send(() => this.#client.GET('/todos'), 'load');
		if (result.ok) this.items = result.data!;
		this.loading = false;
	}

	/** Each change returns whether it succeeded. */
	async add(title: string): Promise<boolean> {
		const result = await this.#send(() => this.#client.POST('/todos', { body: { title } }));
		if (result.ok) this.items = [...this.items, result.data!];
		return result.ok;
	}

	async toggle(todo: Todo): Promise<boolean> {
		return this.#update(todo, { done: !todo.done });
	}

	async rename(todo: Todo, title: string): Promise<boolean> {
		return this.#update(todo, { title });
	}

	async remove(todo: Todo): Promise<boolean> {
		const params = { path: { id: todo.id } };
		const result = await this.#whileBusy(todo, () =>
			this.#send(() => this.#client.DELETE('/todos/{id}', { params }))
		);
		if (result.ok) this.items = this.items.filter((item) => item.id !== todo.id);
		return result.ok;
	}

	dismissError() {
		this.actionError = null;
	}

	async #update(todo: Todo, body: { title?: string; done?: boolean }): Promise<boolean> {
		const params = { path: { id: todo.id } };
		const result = await this.#whileBusy(todo, () =>
			this.#send(() => this.#client.PATCH('/todos/{id}', { params, body }))
		);
		if (result.ok)
			this.items = this.items.map((item) => (item.id === todo.id ? result.data! : item));
		return result.ok;
	}

	/** Marks the todo busy until the request answers (#send never throws). */
	async #whileBusy<T>(todo: Todo, request: () => Promise<T>): Promise<T> {
		this.busy.add(todo.id);
		const result = await request();
		this.busy.delete(todo.id);
		return result;
	}

	/** Sends a request; on failure (API error or no connection) records the error and returns not ok. */
	async #send<T>(
		request: () => Promise<Result<T>>,
		kind: 'load' | 'action' = 'action'
	): Promise<{ ok: true; data?: T } | { ok: false }> {
		let error: ApiError;
		try {
			const result = await request();
			if (!result.error) return { ok: true, data: result.data };
			error = result.error;
		} catch (cause) {
			error = { code: 0, message: String(cause) };
		}
		if (kind === 'load') this.loadError = error;
		else this.actionError = error;
		return { ok: false };
	}
}
