<script lang="ts">
	import './app.css';
	import { Button } from '$lib/components/ui/button/index.js';

	let a = $state('this is a test');

	$effect(() => {
		console.log('a is', a);
	});
</script>

<h1>lurra</h1>
<Button>Click me</Button>
