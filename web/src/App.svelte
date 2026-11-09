<script lang="ts">
	import './app.css';
	import { Router } from 'sv-router';
	import LanguageSelector from '$lib/components/LanguageSelector.svelte';
	import { m } from './paraglide/messages.js';
	import { route } from './router';
</script>

<header class="flex justify-end p-4">
	<LanguageSelector />
</header>
<h1>lurra</h1>
<p>{m.current_route({ path: route.pathname })}</p>
<Router />
