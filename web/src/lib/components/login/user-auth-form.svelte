<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { toast } from 'svelte-sonner';
	import { cn } from '$lib/utils';
	import { i18n } from '$lib/i18n';
	import { getCallback } from '$lib/callback';
	import { signIn } from '$lib/auth-client';
	import { writable } from 'svelte/store';

	const email = writable('');
	const password = writable('');
	const loading = writable(false);

	async function onSubmit() {
		loading.set(true);

		await signIn.email(
			{
				email: $email,
				password: $password,
				callbackURL: i18n.resolveRoute(getCallback())
			},
			{
				async onSuccess(context) {
					// TODO: welcome message
					toast.success('Logged in successfully!');
				},
				onError(context) {
					toast.error(context.error.message);
				}
			}
		);

		loading.set(false);
	}

	function updateOAuth2MapLoading(provider: string, loading: boolean) {
		const oauth2Provider = oauth2Map.get(provider);
		if (oauth2Provider) {
			oauth2Provider.loading = loading;
		}
		oauth2Map = oauth2Map;
	}

	async function authWithOAuth2(provider: string) {
		updateOAuth2MapLoading(provider, true);

		await signIn.social(
			{
				provider: provider as
					| 'apple'
					| 'discord'
					| 'facebook'
					| 'github'
					| 'google'
					| 'microsoft'
					| 'spotify'
					| 'twitch'
					| 'twitter'
					| 'dropbox'
					| 'linkedin'
					| 'gitlab'
					| 'tiktok'
					| 'reddit'
					| 'roblox'
					| 'vk'
					| 'kick',
				callbackURL: i18n.resolveRoute(getCallback()) + '/123'
			},
			{
				onSuccess(data) {
					// TODO: welcome message
					toast.success('Logged in successfully!');
				},
				onError(context) {
					toast.error(context.error.message);
				}
			}
		);

		updateOAuth2MapLoading(provider, false);
	}

	interface OAuth2 {
		name: string;
		icon: string;
		enabled?: boolean;
		loading: boolean;
	}

	let oauth2Map = new Map<string, OAuth2>([
		[
			'apple',
			{
				name: 'Apple',
				icon: 'ph:apple-logo',
				loading: false
			}
		],
		[
			'facebook',
			{
				name: 'Facebook',
				icon: 'ph:facebook-logo',
				loading: false
			}
		],
		[
			'github',
			{
				name: 'GitHub',
				icon: 'ph:github-logo',
				loading: false
			}
		],
		[
			'google',
			{
				name: 'Google',
				icon: 'ph:google-logo',
				loading: false
			}
		],
		[
			'twitter',
			{
				name: 'X',
				icon: 'ph:x-logo',
				loading: false
			}
		]
	]);

	function splitMap<K, V>(map: Map<K, V>): [Map<K, V>, Map<K, V>] {
		const map1 = new Map<K, V>();
		const map2 = new Map<K, V>();
		const entries = Array.from(map.entries());

		for (let i = 0; i < entries.length; i++) {
			if (i < 5) {
				map1.set(entries[i][0], entries[i][1]);
			} else {
				map2.set(entries[i][0], entries[i][1]);
			}
		}

		return [map1, map2];
	}
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
					disabled={$loading}
					bind:value={$email}
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
					disabled={$loading}
					bind:value={$password}
				/>
			</div>
			<Button type="submit" disabled={$loading} class="[&_svg]:size-5">
				{#if $loading}
					<Icon icon="ph:spinner-gap" class="animate-spin" />
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
	<div class="flex-col space-y-2">
		{#each splitMap(oauth2Map) as oauthRow}
			<div class="flex justify-evenly space-x-2">
				{#each oauthRow as [provider, oauth2]}
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger
								class={cn(
									buttonVariants({ variant: 'outline' }),
									'disabled:pointer-events-auto disabled:cursor-not-allowed [&_svg]:size-5'
								)}
								disabled={oauth2.loading}
								onclick={() => authWithOAuth2(provider)}
							>
								{#if oauth2.loading}
									<Icon icon="ph:spinner-gap" class="animate-spin" />
								{:else}
									<Icon icon={oauth2.icon} class="strike" />
								{/if}
							</Tooltip.Trigger>
							<Tooltip.Content>
								<p>{oauth2.name}</p>
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				{/each}
			</div>
		{/each}
	</div>
</div>
