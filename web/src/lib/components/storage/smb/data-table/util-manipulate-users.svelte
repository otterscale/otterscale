<script lang="ts" module>
	import Icon from '@iconify/svelte';

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
		users = $bindable(),
		type,
		invalid = $bindable(),
		required = $bindable()
	}: {
		users: SMBShare_SecurityConfig_User[];
		type: 'create' | 'update';
		invalid?: boolean;
		required?: boolean;
	} = $props();

	if (!users) {
		users = [];
	}

	const defaults = {} as SMBShare_SecurityConfig_User;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	$effect(() => {
		invalid = !(users && users.length > 0);
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger
		variant="default"
		class={cn(
			'w-full ring-1 ring-primary',
			required && !(users && users.length > 0)
				? 'text-destructive ring-destructive'
				: buttonVariants({ variant: 'outline' })
		)}
	>
		{#if type === 'create'}
			{m.create_users()}
		{:else if type === 'update'}
			{m.edit_users()}
		{/if}
	</Modal.Trigger>
	<Modal.Content>
		{#if users.length > 0}
			<div class="max-h-40 overflow-y-auto rounded-lg border p-2">
				{#each users as user, index (index)}
					<div class="flex items-center gap-2 rounded-lg p-2">
						<div class={cn('flex size-8 items-center justify-center rounded-full border-2')}>
							<Icon icon="ph:user" class="size-5" />
						</div>

						<div class="flex flex-col gap-1">
							<p class="text-xs text-muted-foreground">{m.user()}</p>
							<p class="text-sm">{user.username}</p>
						</div>

						<div class="ml-auto">
							<Button
								variant="ghost"
								size="icon"
								onclick={() => {
									users.splice(index, 1);
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
		<SingleInput.General type="text" bind:value={request.username} required={users.length === 0} />

		<Form.Label>{m.password()}</Form.Label>
		<SingleInput.General
			type="password"
			bind:value={request.password}
			required={users.length === 0}
		/>

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
					disabled={invalid}
					onclick={() => {
						users = [];
					}}
				>
					{m.clear()}
				</Button>
				<Button
					onclick={() => {
						users = [...users, request];
					}}
				>
					{m.add()}
				</Button>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
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
