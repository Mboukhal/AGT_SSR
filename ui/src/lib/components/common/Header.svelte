<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { MoonIcon, SunIcon, KeyRound } from '@lucide/svelte';
	import { toggleMode } from 'mode-watcher';
	import { Button } from '../ui/button';

	class SessionStore {
		private session: { email: string; picture: string } | null = null;

		constructor() {
			this.fetchSession();
		}

		async fetchSession() {
			try {
				const response = await fetch('/api/v1/auth/session');
				if (response.ok) {
					this.session = await response.json();
				} else {
					this.session = null;
				}
			} catch (error) {
				console.error('Failed to fetch session:', error);
				this.session = null;
			}
		}

		get() {
			return this.session;
		}

		isAuthenticated() {
			return this.session !== null;
		}
	}

	const sessionStore = new SessionStore();

	const HEADER_SIZEZ = 50;
</script>

<header
	class="fixed inset-0 z-50 flex items-center justify-between gap-2 bg-background/50 px-6 shadow-sm backdrop-blur-sm ease-linear"
	style="height: {HEADER_SIZEZ}px;"
>
	<a href="/" class="flex items-center gap-2">
		<h1 class="0 font-bold text-primary brightness-200">{import.meta.env.APP_NAME}</h1>
	</a>

	<div class="flex items-center gap-4">
		<Button
			onclick={toggleMode}
			variant="secondary"
			class="rounded-full px-3 transition-all! duration-500 "
		>
			<SunIcon class="hidden size-3 dark:block" />
			<MoonIcon class="size-3 dark:hidden" />
		</Button>
		{#if sessionStore.isAuthenticated()}
			<Avatar.Root class="cursor-pointer">
				<Avatar.Image class="h-8 w-8" src={sessionStore.get()?.picture} />
				<Avatar.Fallback class="h-8 w-8">
					{sessionStore.get()?.email.charAt(0).toUpperCase()}
				</Avatar.Fallback>
			</Avatar.Root>
		{:else}
			<Button variant="outline" class="rounded-full px-3" href="/login">
				<KeyRound class="size-3" />
			</Button>
		{/if}

		<!-- <DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-fir" align="end">
				<DropdownMenu.Label>
					{sessionStore.get()?.email || import.meta.env.APP_NAME}
				</DropdownMenu.Label>
				<DropdownMenu.Group></DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={toggleMode}>
					<SunIcon class="hidden size-4 dark:block" />
					<MoonIcon class="size-4 dark:hidden" />
					<span>Change Theme</span>
				</DropdownMenu.Item>
				{#if sessionStore.isAuthenticated()}
					<a href="/api/v1/auth/logout">
						<DropdownMenu.Item class="text-destructive!">
							<Cable class=" h-4 w-4" />
							<span>Sign Out</span>
						</DropdownMenu.Item>
					</a>
				{:else}
					<a href="/login">
						<DropdownMenu.Item>
							<Cable class=" h-4 w-4" />
							<span>Sign In</span>
						</DropdownMenu.Item>
					</a>
				{/if}
			</DropdownMenu.Content>
		</DropdownMenu.Root> -->
	</div>
</header>

<div style="height: {HEADER_SIZEZ}px;"></div>
