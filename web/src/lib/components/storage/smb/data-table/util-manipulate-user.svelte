<script lang="ts" module>
	import type { SMBShare_SecurityConfig_User } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		value = $bindable(),
		invalid = $bindable(),
		required = $bindable()
	}: {
		value?: SMBShare_SecurityConfig_User;
		invalid?: boolean;
		required?: boolean;
	} = $props();

	const defaultUser = value ?? ({} as SMBShare_SecurityConfig_User);
	let requestUser = $state({ ...defaultUser });
	function reset() {
		requestUser = { ...defaultUser };
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	$effect(() => {
		invalid = !(value && value.username && value.password);
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger
		variant="default"
		class={cn(
			'w-full ring-1 ring-primary hover:ring-primary',
			required && !(value && value.username && value.password)
				? 'text-destructive ring-destructive'
				: buttonVariants({ variant: 'outline' })
		)}
	>
		{#if !(value && value.username && value.password)}
			{m.create()}
		{:else}
			{m.edit()}
		{/if}
	</Modal.Trigger>
	<Modal.Content>
		<Form.Label>{m.name()}</Form.Label>
		<SingleInput.General
			type="text"
			bind:value={requestUser.username}
			required={!(requestUser && requestUser.username)}
		/>

		<Form.Label>{m.password()}</Form.Label>
		<SingleInput.General
			type="password"
			bind:value={requestUser.password}
			required={!(requestUser && requestUser.password)}
		/>

		<Button
			onclick={() => {
				value = { ...requestUser };
				reset();
				close();
			}}
		>
			{m.confirm()}
		</Button>
	</Modal.Content>
</Modal.Root>
