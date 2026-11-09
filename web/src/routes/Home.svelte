<!-- The home page: `/` and, for now, every other path (router.ts). -->
<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
</script>

<h1 class="mb-4 text-2xl font-bold">{m.home_title()}</h1>
<p>{m.example_message({ username: 'ada' })}</p>
<p>{m.myname({ name: 'Ada Lovelace' })}</p>
