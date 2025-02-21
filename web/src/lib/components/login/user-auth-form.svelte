<script lang="ts">
	import Icon from '@iconify/svelte';

	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils.js';
	import pb from '$lib/pb';
	import { onMount } from 'svelte';

	let className: string | undefined | null = undefined;
	export { className as class };

	let isLoading = false;
	async function onSubmit() {
		isLoading = true;

		setTimeout(() => {
			isLoading = false;
		}, 3000);
	}

	onMount(async () => {
		var authMethods = await pb.collection('users').listAuthMethods();
		authMethods.oauth2.providers.forEach((provider) => {
			updateOAuth2MapEnabled(provider.name, true);
		});
	});

	function updateOAuth2MapEnabled(provider: string, enabled: boolean) {
		const oauth2Provider = oauth2Map.get(provider);
		if (oauth2Provider) {
			oauth2Provider.enabled = enabled;
		}
		oauth2Map = oauth2Map;
	}

	function updateOAuth2MapLoading(provider: string, loading: boolean) {
		const oauth2Provider = oauth2Map.get(provider);
		if (oauth2Provider) {
			oauth2Provider.loading = loading;
		}
		oauth2Map = oauth2Map;
	}

	function callback() {
		const callbackParam = page.url.searchParams.get('callback');
		if (callbackParam) {
			goto(callbackParam);
			return;
		}
		goto('/');
	}

	async function authWithOAuth2(provider: string) {
		updateOAuth2MapLoading(provider, true);
		await pb.collection('users').authWithOAuth2({ provider: provider });
		updateOAuth2MapLoading(provider, false);
		callback();
	}

	interface OAuth2 {
		name: string;
		icon: string;
		enabled?: boolean;
		loading?: boolean;
	}

	let oauth2Map = new Map<string, OAuth2>([
		['apple', { name: 'Apple', icon: 'ph:apple-logo' }],
		['google', { name: 'Google', icon: 'ph:google-logo' }],
		['github', { name: 'GitHub', icon: 'ph:github-logo' }],
		['gitlab', { name: 'GitLab', icon: 'ph:gitlab-logo' }],
		['microsoft', { name: 'Microsoft', icon: 'ph:windows-logo' }]
	]);
</script>

<div class={cn('grid gap-6', className)} {...$$restProps}>
	<form on:submit|preventDefault={onSubmit}>
		<div class="grid gap-2 space-y-2">
			<div class="grid gap-1">
				<Label class="sr-only" for="email">Email</Label>
				<Input
					id="email"
					placeholder="name@example.com"
					type="email"
					autocapitalize="none"
					autocomplete="email"
					autocorrect="off"
					disabled={isLoading}
				/>
			</div>
			<Button type="submit" disabled={isLoading}>
				{#if isLoading}
					<Icon icon="ph:spinner-gap" class="h-5 w-5 animate-spin" />
				{:else}
					<p>Go</p>
				{/if}
			</Button>
		</div>
	</form>
	<div class="relative">
		<div class="absolute inset-0 flex items-center">
			<span class="w-full border-t"></span>
		</div>
		<div class="relative flex justify-center text-xs uppercase">
			<span class="bg-background px-2 text-muted-foreground"> Or continue with </span>
		</div>
	</div>
	<div class="flex justify-evenly space-x-2">
		{#each oauth2Map as [provider, oauth2]}
			<Tooltip.Root>
				<Tooltip.Trigger
					><Button
						variant="outline"
						disabled={!oauth2.enabled}
						on:click={() => authWithOAuth2(provider)}
					>
						{#if oauth2.loading}
							<Icon icon="ph:spinner-gap" class="h-5 w-5 animate-spin" />
						{:else}
							<Icon icon={oauth2.icon} class="strike h-5 w-5" />
						{/if}
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p class={!oauth2.enabled ? 'line-through' : ''}>{oauth2.name}</p>
				</Tooltip.Content>
			</Tooltip.Root>
		{/each}
	</div>
</div>
