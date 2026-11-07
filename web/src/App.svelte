<script lang="ts">
	let a = $state('this is a test');

	$effect(() => {
		console.log('a is', a);
	});
</script>

<h1>lurra</h1>
