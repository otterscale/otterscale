<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';

	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import pb, { Helper, welcomeMessage } from '$lib/pb';
	import { onMount } from 'svelte';
	import { ClientResponseError } from 'pocketbase';
	import { toast } from 'svelte-sonner';
	import { getCallback } from '$lib/utils';

	let email = '';
	let password = '';
	let isLoading = false;

	async function onSubmit() {
		try {
			isLoading = true;
			var m = await pb.collection('users').authWithPassword(email, password);
			if (!m.record.emailVisibility) {
				await welcomeMessage(m.record.id);
				await pb.collection('users').update(m.record.id, { emailVisibility: true });
			}
			goto(getCallback());
		} catch (err) {
			if (err instanceof ClientResponseError) {
				if (!Helper.isEmpty(err.data.data)) {
					toast.error(err.data.data.password.message);
				} else {
					toast.error(err.data.message);
				}
			}
		} finally {
			isLoading = false;
		}
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

	async function authWithOAuth2(provider: string) {
		try {
			updateOAuth2MapLoading(provider, true);
			var m = await pb.collection('users').authWithOAuth2({ provider: provider });
			if (!m.record.emailVisibility) {
				await welcomeMessage(m.record.id);
				await pb.collection('users').update(m.record.id, { emailVisibility: true });
			}
			goto(getCallback());
		} catch {
			toast.error('Authentication failed. Please try again.');
		} finally {
			updateOAuth2MapLoading(provider, false);
		}
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

<div class="grid gap-6" {...$$restProps}>
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
					bind:value={email}
				/>
			</div>
			<div class="grid gap-1">
				<Label class="sr-only" for="password">Password</Label>
				<Input
					id="password"
					placeholder="********"
					type="password"
					autocapitalize="none"
					autocomplete="current-password"
					disabled={isLoading}
					bind:value={password}
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
				<Tooltip.Trigger asChild>
					<Button
						variant="outline"
						class="disabled:pointer-events-auto disabled:cursor-not-allowed"
						disabled={oauth2.loading || !oauth2.enabled}
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
					<p>{oauth2.name}</p>
				</Tooltip.Content>
			</Tooltip.Root>
		{/each}
	</div>
</div>
