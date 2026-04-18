<script lang="ts">
	import './layout.css';
	import { ModeWatcher, setMode } from 'mode-watcher';
	import { onMount } from 'svelte';

	const MODE_STORAGE_KEY = 'mode-watcher-mode';

	onMount(() => {
		const savedMode = localStorage.getItem(MODE_STORAGE_KEY);
		if (savedMode !== 'light' && savedMode !== 'dark') {
			setMode('light');
		}
	});
	let { children } = $props();
</script>

<div class="transition-colors! duration-200">
	{@render children()}
	<ModeWatcher modeStorageKey={MODE_STORAGE_KEY} />
</div>
