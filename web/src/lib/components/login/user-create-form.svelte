<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { i18n } from '$lib/i18n';
	import { signUp } from '$lib/auth-client';
	import { writable } from 'svelte/store';
	import { goto } from '$app/navigation';

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
				onSuccess() {
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

<div class="grid gap-6" {...$$restProps}>
	<form on:submit|preventDefault={onSubmit}>
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
</div>
