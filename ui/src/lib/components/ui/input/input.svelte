<script lang="ts" module>
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLInputAttributes } from "svelte/elements";

	export type InputProps = WithElementRef<HTMLInputAttributes>;
</script>

<script lang="ts">
	let {
		class: className,
		type = "text",
		ref = $bindable(null),
		...restProps
	}: InputProps = $props();
</script>

<input
	bind:this={ref}
	{type}
	data-slot="input"
	class={cn(
		"file:text-foreground placeholder:text-muted-foreground dark:bg-input/30 border-input h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:text-sm file:font-medium focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-3 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm dark:file:bg-transparent",
		"aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50",
		className
	)}
	{...restProps}
/>
