<script lang="ts">
	import './app.css';
	import { Router } from 'sv-router';
	import LanguageSelector from '$lib/components/LanguageSelector.svelte';
	import ThemeSelector from '$lib/components/ThemeSelector.svelte';
	import HouseIcon from '@lucide/svelte/icons/house';
	import { m } from '$lib/paraglide/messages.js';
	import { p } from './router';

	// the browser tab title, in the current language (index.html has the English one)
	$effect(() => {
		document.title = m.app_title();
	});
</script>

<!-- for keyboard users: hidden until focused with Tab, skips past the header -->
<a
	href="#main"
	class="btn sr-only focus:not-sr-only focus:absolute focus:start-2 focus:top-2 focus:z-50"
>
	{m.skip_to_content()}
</a>

<!-- a full-height column, so the footer stays at the bottom on short pages -->
<div class="flex min-h-screen flex-col">
	<header class="navbar bg-base-200">
		<div class="navbar-start">
			<a href={p('/')} class="btn btn-square btn-ghost" aria-label={m.nav_home()}>
				<HouseIcon class="size-5" />
			</a>
		</div>
		<!-- settings, not navigation: no <nav> until there are links to pages -->
		<div class="navbar-end gap-2">
			<ThemeSelector />
			<LanguageSelector />
		</div>
	</header>

	<!-- each page renders its own <h1> -->
	<main id="main" class="mx-auto w-full max-w-3xl flex-1 p-4">
		<Router />
	</main>

	<footer class="footer footer-center bg-base-200 p-4 text-sm text-base-content/70">
		<!-- the link ends the sentence in every language, so it needs no rich-text messages -->
		<p>
			{m.footer_note()}
			<a href="https://github.com/enolgor/lurra" class="link">{m.footer_link()}</a>
		</p>
	</footer>
</div>
