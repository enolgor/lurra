<!-- Any path without its own route (the catch-all in router.ts). -->
<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import { p } from '$lib/router';
</script>

<h1 class="mb-4 text-2xl font-bold">{m.not_found_title()}</h1>
<p class="mb-4 text-base-content/70">{m.not_found_text()}</p>
<a href={p('/')} class="link">{m.not_found_home()}</a>
