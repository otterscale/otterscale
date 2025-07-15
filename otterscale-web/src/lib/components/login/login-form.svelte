<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
	import Icon from '@iconify/svelte';
	import { authClient } from '$lib/auth-client';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { privacyPolicyPath, termsOfServicePath } from '$lib/path';
	import Placeholder from '$lib/static/placeholder.svg';

	const { data } = $props();
	const id = $props.id();

	// State management
	const isSignUp = writable(false);
	const loading = writable(false);

	// Form data
	const signInForm = writable({
		email: '',
		password: ''
	});

	const signUpForm = writable({
		firstName: '',
		lastName: '',
		email: '',
		password: ''
	});

	// Social providers configuration
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

	// Helper functions
	function showError(context: any) {
		const message = context.error.message || 'An error occurred. Please try again.';
		toast.error(message);
	}

	function showSuccess(message: string) {
		toast.success(message);
	}

	async function setLoadingState(asyncFn: () => Promise<void>) {
		loading.set(true);
		try {
			await asyncFn();
		} finally {
			loading.set(false);
		}
	}

	// Auth handlers
	async function handleEmailSignUp() {
		console.log($signUpForm.email);
		await setLoadingState(async () => {
			await authClient.signUp.email({
				email: $signUpForm.email,
				password: $signUpForm.password,
				callbackURL: data.nextPath,
				name: `${$signUpForm.firstName} ${$signUpForm.lastName}`,
				fetchOptions: {
					onSuccess: () => showSuccess('Account created successfully!'),
					onError: showError
				}
			});
		});
	}

	async function handleEmailSignIn() {
		await setLoadingState(async () => {
			await authClient.signIn.email({
				email: $signInForm.email,
				password: $signInForm.password,
				callbackURL: data.nextPath,
				fetchOptions: {
					onSuccess: () => showSuccess('Signed in successfully!'),
					onError: showError
				}
			});
		});
	}

	async function handleSocialSignIn(provider: string) {
		await setLoadingState(async () => {
			await authClient.signIn.social({
				provider: provider as any,
				callbackURL: data.nextPath,
				fetchOptions: {
					onSuccess: () => showSuccess('Signed in successfully!'),
					onError: showError
				}
			});
		});
	}

	function toggleMode() {
		isSignUp.update((current) => !current);
	}
</script>

<div class="flex flex-col gap-6">
	<Card.Root class="overflow-hidden p-0">
		<Card.Content class="grid p-0 md:grid-cols-2">
			<!-- Sign In Form -->
			{#if !$isSignUp}
				<form class="p-6 md:p-8" onsubmit={handleEmailSignIn}>
					<div class="flex flex-col gap-6">
						<div class="flex flex-col items-center text-center">
							<h1 class="text-2xl font-bold">Welcome back</h1>
							<p class="text-muted-foreground text-balance">Login to access your account</p>
						</div>

						<div class="grid gap-3">
							<Label for="signin-email-{id}">Email</Label>
							<Input
								id="signin-email-{id}"
								type="email"
								placeholder="name@example.com"
								required
								disabled={$loading}
								bind:value={$signInForm.email}
							/>
						</div>

						<div class="grid gap-3">
							<Label for="signin-password-{id}">Password</Label>
							<Input
								id="signin-password-{id}"
								type="password"
								placeholder="********"
								required
								disabled={$loading}
								bind:value={$signInForm.password}
							/>
						</div>

						<Button type="submit" class="w-full" disabled={$loading}>
							{#if $loading}
								<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
							{:else}
								Login
							{/if}
						</Button>

						<!-- Social Sign In -->
						<div
							class="after:border-border relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t"
						>
							<span class="bg-card text-muted-foreground relative z-10 px-2">
								OR CONTINUE WITH
							</span>
						</div>

						<div class="grid grid-cols-3 gap-4">
							{#each socialProviders as provider}
								<Button
									variant="outline"
									type="button"
									class="w-full"
									disabled={!provider.enabled || $loading}
									onclick={() => handleSocialSignIn(provider.id)}
								>
									{#if $loading}
										<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
									{:else}
										<Icon icon={provider.icon} class="size-5" />
										<span class="sr-only">Login with {provider.label}</span>
									{/if}
								</Button>
							{/each}
						</div>

						<Button variant="link" class="text-muted-foreground" onclick={toggleMode}>
							Don't have an account?
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
			{#if $isSignUp}
				<form class="p-6 md:p-8" onsubmit={handleEmailSignUp}>
					<div class="flex flex-col gap-6">
						<div class="flex flex-col items-center text-center">
							<h1 class="text-2xl font-bold">Hello!</h1>
							<p class="text-muted-foreground text-balance">Create your account to get started</p>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div class="grid gap-3">
								<Label for="first-name-{id}">First Name</Label>
								<Input
									id="first-name-{id}"
									placeholder="Paul"
									required
									disabled={$loading}
									bind:value={$signUpForm.firstName}
								/>
							</div>
							<div class="grid gap-3">
								<Label for="last-name-{id}">Last Name</Label>
								<Input
									id="last-name-{id}"
									placeholder="Smith"
									required
									disabled={$loading}
									bind:value={$signUpForm.lastName}
								/>
							</div>
						</div>

						<div class="grid gap-3">
							<Label for="email-{id}">Email</Label>
							<Input
								id="email-{id}"
								type="email"
								placeholder="name@example.com"
								required
								disabled={$loading}
								bind:value={$signUpForm.email}
							/>
						</div>

						<div class="grid gap-3">
							<Label for="password-{id}">Password</Label>
							<Input
								id="password-{id}"
								type="password"
								placeholder="********"
								required
								disabled={$loading}
								bind:value={$signUpForm.password}
							/>
						</div>

						<Button type="submit" class="w-full" disabled={$loading}>
							{#if $loading}
								<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
							{:else}
								Sign Up
							{/if}
						</Button>

						<Button variant="link" class="text-muted-foreground" onclick={toggleMode}>
							Already have an account?
						</Button>
					</div>
				</form>
			{/if}
		</Card.Content>
	</Card.Root>

	<!-- Terms and Privacy -->
	<div
		class="text-muted-foreground *:[a]:hover:text-primary text-center text-xs text-balance *:[a]:underline *:[a]:underline-offset-4"
	>
		By clicking continue, you agree to our <a href={termsOfServicePath}>Terms of Service</a> and
		<a href={privacyPolicyPath}>Privacy Policy</a>.
	</div>
</div>
