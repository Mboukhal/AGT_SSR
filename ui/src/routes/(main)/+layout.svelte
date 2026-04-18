<script lang="ts">
	import Header from '$lib/components/common/Header.svelte';

	let { children } = $props();
</script>

<main>
	<Header />
	<section class="container mt-8">
		{@render children()}
	</section>
</main>
