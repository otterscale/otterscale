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
	import { signUp } from '$lib/auth-client';
	import { goto, invalidate } from '$app/navigation';

	let { class: className, ...restProps }: HTMLAttributes<HTMLDivElement> = $props();

	const firstName = writable('');
	const lastName = writable('');
	const email = writable('');
	const password = writable('');
	const passwordConfirm = writable('');
	const loading = writable(false);

	async function onSubmit() {
		if ($password !== $passwordConfirm) {
			toast.error('Password confirmation does not match.');
			return;
		}

		loading.set(true);

		await signUp.email({
			email: $email,
			password: $password,
			name: `${$firstName} ${$lastName}`,
			fetchOptions: {
				async onSuccess() {
					await invalidate('app:user');
					toast.success('Account created successfully! Please sign in to continue.');
					goto(i18n.resolveRoute('/'));
				},
				onError(context) {
					toast.error(context.error.message);
				}
			}
		});

		loading.set(false);
	}
</script>

<div class={cn('flex w-full flex-col gap-6', className)} {...restProps}>
	<Card.Root>
		<Card.Header class="text-center">
			<Card.Title class="text-xl">Sign up</Card.Title>
			<Card.Description>Enter your information to create an account</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-6">
			<form onsubmit={onSubmit}>
				<div class="grid gap-2 space-y-2">
					<div class="grid grid-cols-2 gap-4">
						<div class="grid gap-2">
							<Label for="first-name">First name</Label>
							<Input
								id="first-name"
								placeholder="Max"
								autocomplete="given-name"
								disabled={$loading}
								bind:value={$firstName}
								required
							/>
						</div>
						<div class="grid gap-2">
							<Label for="last-name">Last name</Label>
							<Input
								id="last-name"
								placeholder="Robinson"
								autocomplete="family-name"
								disabled={$loading}
								bind:value={$lastName}
								required
							/>
						</div>
					</div>
					<div class="grid gap-2">
						<Label for="email">Email</Label>
						<Input
							id="email"
							placeholder="name@example.com"
							type="email"
							autocapitalize="none"
							autocomplete="email"
							autocorrect="off"
							disabled={$loading}
							bind:value={$email}
							required
						/>
					</div>
					<div class="grid gap-2">
						<Label for="password">Password</Label>
						<Input
							id="password"
							placeholder="********"
							type="password"
							autocapitalize="none"
							autocomplete="current-password"
							disabled={$loading}
							bind:value={$password}
							required
						/>
					</div>
					<div class="grid gap-2">
						<Label for="password-confirm">Confirm password</Label>
						<Input
							id="password-confirm"
							placeholder="********"
							type="password"
							autocapitalize="none"
							autocomplete="current-password"
							disabled={$loading}
							bind:value={$passwordConfirm}
							required
						/>
					</div>
					<Button type="submit" class="[&_svg]:size-5" disabled={$loading}>
						{#if $loading}
							<Icon icon="ph:spinner-gap" class="animate-spin" />
						{:else}
							<p>Create an account</p>
						{/if}
					</Button>
				</div>
			</form>
		</Card.Content>
		<Card.Footer></Card.Footer>
	</Card.Root>
</div>
