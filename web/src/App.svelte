<script lang="ts">
	import './app.css';
	import { Router } from 'sv-router';
	import LanguageSelector from '$lib/components/LanguageSelector.svelte';
	import { route } from './router';
</script>

<header class="flex justify-end p-4">
	<LanguageSelector />
</header>
<h1>lurra</h1>
<!-- a template literal, so Svelte adds no `?? ''` branch that tests could never cover -->
<p>{`you are in route: ${route.pathname}`}</p>
<Router />
