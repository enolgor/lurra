<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import type { Snippet } from 'svelte';

	// One choice in a Dropdown menu: a button that marks the current choice with a check and
	// aria-current, so screen readers announce it too.
	let {
		current,
		onclick,
		children
	}: {
		/** Whether this is the choice in effect. */
		current: boolean;
		onclick: () => void;
		/** The option's content: an icon or flag and its name. */
		children: Snippet;
	} = $props();
</script>

<li>
	<button type="button" aria-current={current} {onclick}>
		{@render children()}
		{#if current}
			<CheckIcon class="ml-auto size-4" />
		{/if}
	</button>
</li>
