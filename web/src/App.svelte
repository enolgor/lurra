<script lang="ts">
	import './app.css';
	import { Router } from 'sv-router';
	import LanguageSelector from '$lib/components/LanguageSelector.svelte';
	import ThemeSelector from '$lib/components/ThemeSelector.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { route } from './router';
</script>

<header class="flex justify-end gap-2 p-4">
	<ThemeSelector />
	<LanguageSelector />
</header>
<h1>lurra</h1>
<p>{m.current_route({ path: route.pathname })}</p>
<Router />
