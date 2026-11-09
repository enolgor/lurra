<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import LanguagesIcon from '@lucide/svelte/icons/languages';
	import { localeFlag, localeName } from '$lib/locale';
	import { m } from '../../paraglide/messages.js';
	import { getLocale, locales, setLocale, type Locale } from '../../paraglide/runtime.js';

	// A native <details> dropdown (daisyUI): it opens and closes on its own, but not on an outside
	// click or Escape, so those are handled below.
	let open = $state(false);
	let dropdown: HTMLDetailsElement;

	// Switches live, without a reload; $lib/i18n.svelte.ts makes setLocale do that.
	function changeLocale(locale: Locale) {
		// typed void | Promise<void> for Paraglide's async strategies; ours is synchronous
		void setLocale(locale);
		open = false;
	}

	function closeOnOutsideClick(event: MouseEvent) {
		if (open && !dropdown.contains(event.target as Node)) open = false;
	}

	function closeOnEscape(event: KeyboardEvent) {
		if (!open || event.key !== 'Escape') return;
		open = false;
		dropdown.querySelector('summary')!.focus();
	}
</script>

<svelte:window onclick={closeOnOutsideClick} onkeydown={closeOnEscape} />

<details class="dropdown dropdown-end" bind:open bind:this={dropdown}>
	<summary class="btn">
		<LanguagesIcon class="size-4" />
		<span>{m.language_selector_trigger()}</span>
	</summary>
	<ul class="menu dropdown-content z-10 w-52 rounded-box bg-base-100 p-2 shadow-sm">
		<li class="menu-title">{m.language_selector_label()}</li>
		{#each locales as locale (locale)}
			<li>
				<button
					type="button"
					aria-current={locale === getLocale()}
					onclick={() => changeLocale(locale)}
				>
					<span aria-hidden="true">{localeFlag(locale)}</span>
					<span lang={locale}>{localeName(locale)}</span>
					{#if locale === getLocale()}
						<CheckIcon class="ml-auto size-4" />
					{/if}
				</button>
			</li>
		{/each}
	</ul>
</details>
