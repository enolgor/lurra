<script lang="ts">
	import type { Snippet } from 'svelte';

	// A daisyUI dropdown on a native <details>: it opens and closes itself on its button, and this
	// adds what <details> lacks. It closes on a click outside, on a click on any button in the menu
	// (every menu item is a choice), and on Escape, which gives focus back to its button.
	let {
		trigger,
		children
	}: {
		/** The button's content. */
		trigger: Snippet;
		/** The menu: <li> items with daisyUI `menu` markup. */
		children: Snippet;
	} = $props();

	let open = $state(false);
	let dropdown: HTMLDetailsElement;
	let button: HTMLElement;

	function closeOnClick(event: MouseEvent) {
		if (!open) return;
		const target = event.target as Element;
		if (!dropdown.contains(target) || target.closest('.dropdown-content button')) open = false;
	}

	function closeOnEscape(event: KeyboardEvent) {
		if (!open || event.key !== 'Escape') return;
		open = false;
		button.focus();
	}
</script>

<svelte:window onclick={closeOnClick} onkeydown={closeOnEscape} />

<details class="dropdown dropdown-end" bind:open bind:this={dropdown}>
	<summary class="btn" bind:this={button}>{@render trigger()}</summary>
	<ul class="menu dropdown-content z-10 w-52 rounded-box bg-base-100 p-2 shadow-sm">
		{@render children()}
	</ul>
</details>
