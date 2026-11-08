<script lang="ts">
	import './app.css';
	import { Router } from 'sv-router';
	import { route } from './router';
</script>

<h1>lurra</h1>
<!-- a template literal, so Svelte adds no `?? ''` branch that tests could never cover -->
<p>{`you are in route: ${route.pathname}`}</p>
<Router />
