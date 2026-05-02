<script lang="ts">
	import Button from "$lib/components/ui/button/button.svelte";
</script>

<div class="w-full max-w-sm mx-auto space-y-4">
	<h1 class="text-2xl font-bold text-center mb-8">Sign In</h1>
	<Button href="/api/auth/microsoft/login" class="w-full" variant="outline">
		Sign in with Microsoft
	</Button>
</div>
