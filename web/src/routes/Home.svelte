<!-- The home page: the todo list, at `/` (router.ts). -->
<script lang="ts">
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import TrashIcon from '@lucide/svelte/icons/trash';
	import XIcon from '@lucide/svelte/icons/x';
	import { onMount } from 'svelte';
	import { searchParams } from 'sv-router';
	import ApiErrorAlert from '$lib/components/ApiErrorAlert.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { filters, TodoList, type Filter, type Todo } from '$lib/todos.svelte';

	/** Tests pass a TodoList with a stubbed API; the app uses the default one. */
	let { todos = new TodoList() }: { todos?: TodoList } = $props();
	onMount(() => {
		void todos.load();
	});

	// The filter lives in the URL (?filter=done), so a reload or a shared link keeps it.
	const isFilter = (value: unknown): value is Filter => filters.includes(value as Filter);
	const filter = $derived.by(() => {
		const value = searchParams.get('filter');
		return isFilter(value) ? value : 'all';
	});
	$effect(() => {
		todos.filter = filter;
	});
	function setFilter(value: Filter) {
		// replace: switching the view doesn't add a history entry
		if (value === 'all') searchParams.delete('filter', undefined, { replace: true });
		else searchParams.set('filter', value, { replace: true });
	}
	const filterNames = { all: m.filter_all, done: m.filter_done, undone: m.filter_undone };

	let newTitle = $state('');
	async function add(event: SubmitEvent) {
		event.preventDefault();
		if (await todos.add(newTitle.trim())) newTitle = '';
	}

	// Inline editing: one todo at a time; Enter or leaving the field saves, Escape cancels, and a
	// blank or unchanged title just closes the field.
	let editing = $state<number | null>(null);
	function save(todo: Todo, title: string) {
		editing = null;
		title = title.trim();
		if (title && title !== todo.title) void todos.rename(todo, title);
	}
	function onEditKey(event: KeyboardEvent, todo: Todo) {
		if (event.key === 'Enter') save(todo, (event.currentTarget as HTMLInputElement).value);
		if (event.key === 'Escape') editing = null;
	}
	function toggle(event: MouseEvent, todo: Todo) {
		// the box shows the saved state: it changes when the API answers, not on the click
		event.preventDefault();
		void todos.toggle(todo);
	}
</script>

<h1 class="mb-4 text-2xl font-bold">{m.home_title()}</h1>

<form class="mb-4 flex w-full gap-2" onsubmit={add}>
	<input
		class="input w-full"
		bind:value={newTitle}
		placeholder={m.todo_new_placeholder()}
		aria-label={m.todo_new_label()}
	/>
	<button class="btn btn-primary" disabled={!newTitle.trim()}>{m.todo_add()}</button>
</form>

<div role="tablist" aria-label={m.filter_label()} class="tabs tabs-border mb-2">
	{#each filters as value (value)}
		<button
			type="button"
			role="tab"
			class="tab"
			class:tab-active={value === filter}
			aria-selected={value === filter}
			onclick={() => setFilter(value)}
		>
			{filterNames[value]()}
		</button>
	{/each}
</div>

{#if todos.loading}
	<div role="status" class="flex justify-center py-8">
		<span class="loading loading-lg loading-spinner" aria-hidden="true"></span>
		<span class="sr-only">{m.loading()}</span>
	</div>
{:else if todos.loadError}
	<ApiErrorAlert error={todos.loadError}>
		{#snippet action()}
			<button type="button" class="btn btn-sm" onclick={() => todos.load()}>{m.try_again()}</button>
		{/snippet}
	</ApiErrorAlert>
{:else if todos.items.length === 0}
	<p class="py-8 text-center text-base-content/60">{m.todo_empty()}</p>
{:else if todos.visible.length === 0}
	<p class="py-8 text-center text-base-content/60">{m.todo_none_filtered()}</p>
{:else}
	<ul class="list rounded-box bg-base-100">
		{#each todos.visible as todo (todo.id)}
			{@const busy = todos.busy.has(todo.id)}
			<li class="list-row items-center">
				<input
					type="checkbox"
					class="checkbox"
					checked={todo.done}
					disabled={busy}
					aria-label={todo.title}
					onclick={(event) => toggle(event, todo)}
				/>
				{#if editing === todo.id}
					<input
						class="input input-sm"
						value={todo.title}
						aria-label={m.todo_title_label()}
						onkeydown={(event) => onEditKey(event, todo)}
						onblur={(event) => save(todo, event.currentTarget.value)}
						{@attach (input: HTMLInputElement) => input.focus()}
					/>
				{:else}
					<span class={todo.done ? 'text-base-content/60 line-through' : ''}>{todo.title}</span>
				{/if}
				<div class="flex gap-1">
					<button
						type="button"
						class="btn btn-square btn-ghost btn-sm"
						disabled={busy}
						aria-label={m.todo_edit({ title: todo.title })}
						onclick={() => (editing = todo.id)}
					>
						<PencilIcon class="size-4" />
					</button>
					<button
						type="button"
						class="btn btn-square btn-ghost btn-sm"
						disabled={busy}
						aria-label={m.todo_delete({ title: todo.title })}
						onclick={() => todos.remove(todo)}
					>
						<TrashIcon class="size-4" />
					</button>
				</div>
			</li>
		{/each}
	</ul>
{/if}

{#if todos.actionError}
	<div class="toast toast-end">
		<ApiErrorAlert error={todos.actionError}>
			{#snippet action()}
				<button
					type="button"
					class="btn btn-square btn-ghost btn-sm"
					aria-label={m.close()}
					onclick={() => todos.dismissError()}
				>
					<XIcon class="size-4" />
				</button>
			{/snippet}
		</ApiErrorAlert>
	</div>
{/if}
