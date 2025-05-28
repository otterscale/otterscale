<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { i18n } from '$lib/i18n';
	import { cn } from '$lib/utils';
	import { appendCallback, getCallback } from '$lib/callback';
	import { signIn } from '$lib/auth-client';
	import { invalidate } from '$app/navigation';

	let { class: className, ...restProps }: HTMLAttributes<HTMLDivElement> = $props();

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

	async function authWithOAuth2(oauth2: OAuth2) {
		oauth2.loading = true;

		await signIn.social(
			{
				provider: oauth2.provider,
				callbackURL: i18n.resolveRoute(getCallback())
			},
			{
				async onSuccess() {
					// TODO: welcome message
					await invalidate('app:user');
					toast.success('Logged in successfully!');
				},
				onError(context) {
					toast.error(context.error.message);
					oauth2.loading = false;
				}
			}
		);
	}

	interface OAuth2 {
		provider:
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
			| 'kick';
		name: string;
		icon: string;
		loading: boolean;
	}

	let oauth2List = $state<OAuth2[]>([
		{
			provider: 'apple',
			name: 'Apple',
			icon: 'ph:apple-logo',
			loading: false
		},
		{
			provider: 'facebook',
			name: 'Facebook',
			icon: 'ph:facebook-logo',
			loading: false
		},
		{
			provider: 'github',
			name: 'GitHub',
			icon: 'ph:github-logo',
			loading: false
		},
		{
			provider: 'google',
			name: 'Google',
			icon: 'ph:google-logo',
			loading: false
		},
		{
			provider: 'twitter',
			name: 'X',
			icon: 'ph:x-logo',
			loading: false
		}
	]);
</script>

<div class={cn('flex w-full flex-col gap-6', className)} {...restProps}>
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">Login</Card.Title>
			<Card.Description>Enter your email to access your account</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-6">
			<form onsubmit={onSubmit}>
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
				<div class="flex justify-evenly space-x-2">
					{#each oauth2List as oauth2}
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger
									class={cn(
										buttonVariants({ variant: 'outline' }),
										'disabled:pointer-events-auto disabled:cursor-not-allowed [&_svg]:size-5'
									)}
									disabled={oauth2.loading}
									onclick={() => authWithOAuth2(oauth2)}
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
			</div>
		</Card.Content>
	</Card.Root>
</div>
