<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { authClient } from '$lib/auth-client';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { m } from '$lib/paraglide/messages';
	import Placeholder from '$lib/static/placeholder.svg';

	const { data } = $props();
	const id = $props.id();

	// Types
	interface LoadingState {
		email: boolean;
		apple: boolean;
		github: boolean;
		google: boolean;
	}

	interface SignInForm {
		email: string;
		password: string;
	}

	interface SignUpForm {
		firstName: string;
		lastName: string;
		email: string;
		password: string;
	}

	// State
	const signUp = writable(false);
	const loading = writable<LoadingState>({
		email: false,
		apple: false,
		github: false,
		google: false
	});

	const signInForm = writable<SignInForm>({
		email: '',
		password: ''
	});

	const signUpForm = writable<SignUpForm>({
		firstName: '',
		lastName: '',
		email: '',
		password: ''
	});

	// Constants
	const socialProviders = [
		{ id: 'apple', icon: 'streamline-logos:apple-logo-solid', label: 'Apple', enabled: data.apple },
		{
			id: 'github',
			icon: 'streamline-logos:github-logo-2-solid',
			label: 'GitHub',
			enabled: data.github
		},
		{
			id: 'google',
			icon: 'streamline-logos:google-logo-solid',
			label: 'Google',
			enabled: data.google
		}
	];

	// Utility functions
	const showError = (context: any) => {
		const message = context.error.message || 'An error occurred. Please try again.';
		toast.error(message);
	};

	const setLoadingState = async (provider: string, asyncFn: () => Promise<void>) => {
		loading.update((state) => ({ ...state, [provider]: true }));
		try {
			await asyncFn();
		} finally {
			await new Promise((resolve) => setTimeout(resolve, 1000));
			loading.update((state) => ({ ...state, [provider]: false }));
		}
	};

	// Auth handlers
	const handleEmailSignUp = async () => {
		await setLoadingState('email', async () => {
			await authClient.signUp.email({
				email: $signUpForm.email,
				password: $signUpForm.password,
				name: `${$signUpForm.firstName} ${$signUpForm.lastName}`,
				fetchOptions: {
					onSuccess: () => goto(data.nextPath),
					onError: showError
				}
			});
		});
	};

	const handleEmailSignIn = async () => {
		await setLoadingState('email', async () => {
			await authClient.signIn.email({
				email: $signInForm.email,
				password: $signInForm.password,
				callbackURL: data.nextPath,
				fetchOptions: {
					onError: showError
				}
			});
		});
	};

	const handleSocialSignIn = async (provider: string) => {
		await setLoadingState(provider, async () => {
			await authClient.signIn.social({
				provider: provider as any,
				callbackURL: data.nextPath,
				fetchOptions: {
					onError: showError
				}
			});
		});
	};

	const toggleMode = () => signUp.update((current) => !current);

	// Form handlers
	const handleSignInSubmit = (e: Event) => {
		e.preventDefault();
		handleEmailSignIn();
	};

	const handleSignUpSubmit = (e: Event) => {
		e.preventDefault();
		handleEmailSignUp();
	};
</script>

<Card.Root class="overflow-hidden p-0">
	<Card.Content class="grid p-0 md:grid-cols-2">
		<!-- Sign In Form -->
		{#if !$signUp}
			<form class="p-6 md:p-8" onsubmit={handleSignInSubmit}>
				<div class="flex flex-col gap-6">
					<div class="flex flex-col items-center text-center">
						<h1 class="text-2xl font-bold">{m.login_title()}</h1>
						<p class="text-muted-foreground text-balance">{m.login_description()}</p>
					</div>

					<div class="grid gap-3">
						<Label for="signin-email-{id}">{m.email()}</Label>
						<Input
							id="signin-email-{id}"
							type="email"
							placeholder="name@example.com"
							required
							disabled={$loading.email}
							bind:value={$signInForm.email}
						/>
					</div>

					<div class="grid gap-3">
						<Label for="signin-password-{id}">{m.password()}</Label>
						<Input
							id="signin-password-{id}"
							type="password"
							placeholder="********"
							required
							disabled={$loading.email}
							bind:value={$signInForm.password}
						/>
					</div>

					<Button type="submit" class="w-full" disabled={$loading.email}>
						{#if $loading.email}
							<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
						{:else}
							{m.login()}
						{/if}
					</Button>

					<!-- Social Sign In -->
					<div
						class="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t"
					>
						<span class="bg-card text-muted-foreground relative z-10 px-2">{m.login_divider()}</span
						>
					</div>

					<div class="grid grid-cols-3 gap-4">
						{#each socialProviders as provider}
							<Button
								variant="outline"
								type="button"
								class="w-full"
								disabled={!provider.enabled || $loading[provider.id as keyof LoadingState]}
								onclick={() => handleSocialSignIn(provider.id)}
							>
								{#if $loading[provider.id as keyof LoadingState]}
									<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
								{:else}
									<Icon icon={provider.icon} class="size-5" />
									<span class="sr-only">Login with {provider.label}</span>
								{/if}
							</Button>
						{/each}
					</div>

					<Button variant="link" class="text-muted-foreground" onclick={toggleMode}>
						{m.login_toggle()}
					</Button>
				</div>
			</form>
		{/if}

		<!-- Placeholder Image -->
		<div class="bg-muted relative hidden md:block">
			<img
				src={Placeholder}
				alt="placeholder"
				class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
			/>
		</div>

		<!-- Sign Up Form -->
		{#if $signUp}
			<form class="p-6 md:p-8" onsubmit={handleSignUpSubmit}>
				<div class="flex flex-col gap-6">
					<div class="flex flex-col items-center text-center">
						<h1 class="text-2xl font-bold">{m.sign_up_title()}</h1>
						<p class="text-muted-foreground text-balance">{m.sign_up_description()}</p>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div class="grid gap-3">
							<Label for="first-name-{id}">{m.first_name()}</Label>
							<Input
								id="first-name-{id}"
								placeholder="Paul"
								required
								disabled={$loading.email}
								bind:value={$signUpForm.firstName}
							/>
						</div>
						<div class="grid gap-3">
							<Label for="last-name-{id}">{m.last_name()}</Label>
							<Input
								id="last-name-{id}"
								placeholder="Smith"
								required
								disabled={$loading.email}
								bind:value={$signUpForm.lastName}
							/>
						</div>
					</div>

					<div class="grid gap-3">
						<Label for="email-{id}">{m.email()}</Label>
						<Input
							id="email-{id}"
							type="email"
							placeholder="name@example.com"
							required
							disabled={$loading.email}
							bind:value={$signUpForm.email}
						/>
					</div>

					<div class="grid gap-3">
						<Label for="password-{id}">{m.password()}</Label>
						<Input
							id="password-{id}"
							type="password"
							placeholder="********"
							required
							disabled={$loading.email}
							bind:value={$signUpForm.password}
						/>
					</div>

					<Button type="submit" class="w-full" disabled={$loading.email}>
						{#if $loading.email}
							<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
						{:else}
							{m.sign_up()}
						{/if}
					</Button>

					<Button variant="link" class="text-muted-foreground" onclick={toggleMode}>
						{m.sign_up_toggle()}
					</Button>
				</div>
			</form>
		{/if}
	</Card.Content>
</Card.Root>
