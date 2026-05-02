<script lang="ts">
	import { auth } from "$lib/stores/auth.svelte";
	import Button from "$lib/components/ui/button/button.svelte";

	auth.init();
</script>

<main class="min-h-screen flex items-center justify-center">
	<div class="text-center space-y-4">
		<h1 class="text-2xl font-bold">Welcome, {auth.user?.username}!</h1>
		<div class="flex gap-2 justify-center">
			<Button href="/about">About</Button>
			<Button href="/api/auth/logout" variant="outline">Logout</Button>
		</div>
	</div>
</main>
