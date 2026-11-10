<script lang="ts">
	import LanguagesIcon from '@lucide/svelte/icons/languages';
	import Dropdown from '$lib/components/Dropdown.svelte';
	import MenuOption from '$lib/components/MenuOption.svelte';
	import { localeFlag, localeName } from '$lib/locale';
	import { m } from '$lib/paraglide/messages.js';
	import { getLocale, locales, setLocale, type Locale } from '$lib/paraglide/runtime.js';

	// Switches live, without a reload; $lib/i18n.svelte.ts makes setLocale do that.
	function changeLocale(locale: Locale) {
		// typed void | Promise<void> for Paraglide's async strategies; ours is synchronous
		void setLocale(locale);
	}
</script>

<Dropdown>
	{#snippet trigger()}
		<LanguagesIcon class="size-4" />
		<span>{m.language_selector_trigger()}</span>
	{/snippet}
	<li class="menu-title">{m.language_selector_label()}</li>
	{#each locales as locale (locale)}
		<MenuOption current={locale === getLocale()} onclick={() => changeLocale(locale)}>
			<span aria-hidden="true">{localeFlag(locale)}</span>
			<span lang={locale}>{localeName(locale)}</span>
		</MenuOption>
	{/each}
</Dropdown>
