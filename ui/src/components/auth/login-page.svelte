<!-- //@ts-ignore -->

<script lang="ts">
	import Else from "./else.svelte";
	import If from "./if.svelte";
	import LoginForm from "./login-form.svelte";
</script>

<div class="w-full max-w-sm mx-auto">
	<h1 class="text-2xl font-bold text-center mb-8">
		{"{{.Name}}"}
		<If condition={import.meta.env.DEV}>
			<span class="text-sm text-gray-500">(Development Mode)</span>
			<Else />
			<span class="text-sm text-gray-500">(Production Mode)</span>
		</If>
	</h1>
	<LoginForm />
</div>
