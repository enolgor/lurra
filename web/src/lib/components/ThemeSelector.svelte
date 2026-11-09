<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import MonitorIcon from '@lucide/svelte/icons/monitor';
	import MoonIcon from '@lucide/svelte/icons/moon';
	import SunIcon from '@lucide/svelte/icons/sun';
	import Dropdown from '$lib/components/Dropdown.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { getTheme, setTheme, themes } from '$lib/theme.svelte';

	const icons = { light: SunIcon, dark: MoonIcon, system: MonitorIcon };
	// functions, not strings: called in markup, so they re-render when the language changes
	const names = { light: m.theme_light, dark: m.theme_dark, system: m.theme_system };

	const CurrentIcon = $derived(icons[getTheme()]);
</script>

<Dropdown>
	{#snippet trigger()}
		<CurrentIcon class="size-4" />
		<span>{m.theme_selector_label()}</span>
	{/snippet}
	<li class="menu-title">{m.theme_selector_label()}</li>
	{#each themes as theme (theme)}
		{@const Icon = icons[theme]}
		<li>
			<button type="button" aria-current={theme === getTheme()} onclick={() => setTheme(theme)}>
				<Icon class="size-4" />
				<span>{names[theme]()}</span>
				{#if theme === getTheme()}
					<CheckIcon class="ml-auto size-4" />
				{/if}
			</button>
		</li>
	{/each}
</Dropdown>
