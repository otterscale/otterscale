<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		values = $bindable()
	}: {
		values: string[];
	} = $props();

	if (!values) {
		values = [];
	}

	const defaultUser = '';
	let requestUser = $state(defaultUser);
	function resetUser() {
		requestUser = defaultUser;
	}

	const defaultUsers = $derived([...values]);
	let requestUsers = $derived([...defaultUsers]);
	function resetUsers() {
		requestUsers = [...defaultUsers];
	}

	function reset() {
		resetUser();
		resetUsers();
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger
		variant="default"
		class={cn(
			'w-full ring-1 ring-primary hover:ring-primary',
			buttonVariants({ variant: 'outline' })
		)}
	>
		{#if !(values && values.length > 0)}
			{m.edit()}
		{:else}
			{m.add()}
		{/if}
	</Modal.Trigger>
	<Modal.Content>
		{#if requestUsers.length > 0}
			<div class="max-h-40 overflow-y-auto rounded-lg border p-2">
				{#each requestUsers as user, index (index)}
					<div class="flex items-center gap-2 rounded-lg p-2">
						<div class={cn('flex size-8 items-center justify-center rounded-full border-2')}>
							<Icon icon="ph:user" class="size-5" />
						</div>

						<div class="flex flex-col gap-1">
							<p class="text-xs text-muted-foreground">{m.user()}</p>
							<p class="text-sm">{user}</p>
						</div>

						<div class="ml-auto">
							<Button
								variant="ghost"
								size="icon"
								onclick={() => {
									requestUsers.splice(index, 1);
								}}
							>
								<Icon icon="ph:trash" class="size-4 text-destructive" />
							</Button>
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<div
				class="rounded-lg border border-red-300 bg-destructive/10 p-4 text-center text-xs text-destructive"
			>
				{m.create_users_empty_warning()}
			</div>
		{/if}

		<Form.Label>{m.name()}</Form.Label>
		<SingleInput.General type="text" bind:value={requestUser} />

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<div>
				<Button
					variant="destructive"
					disabled={requestUsers.length === 0}
					onclick={() => {
						reset();
					}}
				>
					{m.clear()}
				</Button>
				<Button
					onclick={() => {
						requestUsers = [...requestUsers, requestUser];
						resetUser();
					}}
				>
					{m.add()}
				</Button>
				<Modal.Action
					disabled={requestUsers.length === 0}
					onclick={() => {
						values = requestUsers;
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</div>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
