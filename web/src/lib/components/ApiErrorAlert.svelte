<script lang="ts">
	import CircleExclamationIcon from '@lucide/svelte/icons/circle-exclamation-point';
	import type { Snippet } from 'svelte';
	import { m } from '$lib/paraglide/messages.js';
	import type { ApiError } from '$lib/todos.svelte';

	let {
		error,
		action
	}: {
		error: ApiError;
		/** A button next to the message, such as "Try again" or close. */
		action?: Snippet;
	} = $props();

	// Translated by the error code when there is one (the status alone isn't specific), else by
	// the status. Message functions, called in markup, so they follow the language live.
	const byCode: Record<NonNullable<ApiError['error']>, () => string> = {
		title_invalid: m.error_title_invalid
	};
	const byStatus: Partial<Record<number, () => string>> = {
		0: m.error_network,
		400: m.error_bad_request,
		404: m.error_not_found,
		413: m.error_too_large,
		500: m.error_server
	};

	const message = $derived(
		error.error ? byCode[error.error] : (byStatus[error.code] ?? m.error_unknown)
	);
</script>

<div role="alert" class="alert alert-error">
	<CircleExclamationIcon class="size-5" aria-hidden="true" />
	<div>
		<p>{message()}</p>
		{#if error.request_id}
			<!-- for bug reports: finds the real error in the server logs -->
			<p class="text-xs opacity-70">{m.error_request_id({ id: error.request_id })}</p>
		{/if}
	</div>
	{#if action}
		{@render action()}
	{/if}
</div>
