<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';

	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import pb, { Helper } from '$lib/pb';
	import { ClientResponseError } from 'pocketbase';
	import { toast } from 'svelte-sonner';
	import { getCallback } from '$lib/utils';
	import { i18n } from '$lib/i18n';

	let email = '';
	let password = '';
	let passwordConfirm = '';
	let firstName = '';
	let lastName = '';
	let isLoading = false;

	async function onSubmit() {
		try {
			isLoading = true;
			await pb.collection('users').create({
				email,
				password,
				passwordConfirm,
				name: `${firstName} ${lastName}`
			});
			await pb.collection('users').authWithPassword(email, password);
			goto(i18n.resolveRoute(getCallback()));
		} catch (err) {
			if (err instanceof ClientResponseError) {
				console.error(err.data);
				if (!Helper.isEmpty(err.data.data.passwordConfirm)) {
					if (err.data.data.passwordConfirm.code === 'validation_values_mismatch') {
						toast.error('Password confirmation does not match.');
					} else {
						toast.error(err.data.data.passwordConfirm.message);
					}
				} else if (!Helper.isEmpty(err.data.data.password)) {
					if (err.data.data.password.code === 'validation_min_text_constraint') {
						toast.error('Password must be at least 8 characters.');
					} else {
						toast.error(err.data.data.password.message);
					}
				} else if (!Helper.isEmpty(err.data.data.email)) {
					if (err.data.data.email.code === 'validation_not_unique') {
						toast.error('Email already exists.');
					} else {
						toast.error(err.data.data.email.message);
					}
				} else {
					toast.error(err.data.message);
				}
			}
		} finally {
			isLoading = false;
		}
		return false;
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
						disabled={isLoading}
						bind:value={firstName}
						required
					/>
				</div>
				<div class="grid gap-2">
					<Label for="last-name">Last name</Label>
					<Input
						id="last-name"
						placeholder="Robinson"
						autocomplete="family-name"
						disabled={isLoading}
						bind:value={lastName}
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
					disabled={isLoading}
					bind:value={email}
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
					disabled={isLoading}
					bind:value={password}
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
					disabled={isLoading}
					bind:value={passwordConfirm}
					required
				/>
			</div>
			<Button type="submit" disabled={isLoading}>
				{#if isLoading}
					<Icon icon="ph:spinner-gap" class="h-5 w-5 animate-spin" />
				{:else}
					<p>Create an account</p>
				{/if}
			</Button>
		</div>
	</form>
</div>
