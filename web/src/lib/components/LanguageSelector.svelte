<script lang="ts">
	import LanguagesIcon from '@lucide/svelte/icons/languages';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { localeFlag, localeName } from '$lib/locale';
	import { m } from '../../paraglide/messages.js';
	import { assertIsLocale, getLocale, locales, setLocale } from '../../paraglide/runtime.js';

	// Switches live, without a reload; $lib/i18n.svelte.ts makes setLocale do that.
	function changeLocale(locale: string) {
		// typed void | Promise<void> for Paraglide's async strategies; ours is synchronous
		void setLocale(assertIsLocale(locale));
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger class={buttonVariants({ variant: 'outline' })}>
		<LanguagesIcon />
		<!-- own element: next to the icon, Svelte would add an untestable `?? ''` branch -->
		<span>{m.language_selector_trigger()}</span>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end" class="w-48">
		<DropdownMenu.Group>
			<DropdownMenu.Label>{m.language_selector_label()}</DropdownMenu.Label>
			<DropdownMenu.Separator />
			<DropdownMenu.RadioGroup value={getLocale()} onValueChange={changeLocale}>
				{#each locales as locale (locale)}
					<DropdownMenu.RadioItem value={locale} closeOnSelect>
						<span aria-hidden="true">{localeFlag(locale)}</span>
						<span lang={locale}>{localeName(locale)}</span>
					</DropdownMenu.RadioItem>
				{/each}
			</DropdownMenu.RadioGroup>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
