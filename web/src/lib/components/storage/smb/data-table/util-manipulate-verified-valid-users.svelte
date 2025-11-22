<script lang="ts" module>
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';

	import {
		StorageService,
		type SMBShare_SecurityConfig_User
	} from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
	import { getContext } from 'svelte';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
</script>

<script lang="ts">
	let {
		validUsers = $bindable(),
		type,
		realm,
		joinSource
	}: {
		validUsers: string[];
		type: 'create' | 'update';
		realm: string;
		joinSource?: SMBShare_SecurityConfig_User;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	if (!validUsers) {
		validUsers = [];
	}

	const defaults = '';
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	const isNull = $derived(validUsers.length === 0);
</script>

<Modal.Root bind:open>
	<Modal.Trigger
		variant="default"
		class={cn('w-full ring-1 ring-primary', buttonVariants({ variant: 'outline' }))}
	>
		{#if type === 'create'}
			{m.create_users()}
		{:else if type === 'update'}
			{m.edit_users()}
		{/if}
	</Modal.Trigger>
	<Modal.Content>
		{#if validUsers.length > 0}
			<div class="max-h-40 overflow-y-auto rounded-lg border p-2">
				{#each validUsers as validUser, index (index)}
					<div class="flex items-center gap-2 rounded-lg p-2">
						<div class={cn('flex size-8 items-center justify-center rounded-full border-2')}>
							<Icon icon="ph:user" class="size-5" />
						</div>

						<div class="flex flex-col gap-1">
							<p class="text-xs text-muted-foreground">{m.user()}</p>
							<p class="text-sm">{validUser}</p>
						</div>

						<div class="ml-auto">
							<Button
								variant="ghost"
								size="icon"
								onclick={() => {
									validUsers.splice(index, 1);
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
		<SingleInput.General type="text" bind:value={request} />

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
					disabled={isNull}
					onclick={() => {
						validUsers = [];
					}}
				>
					{m.clear()}
				</Button>
				<Button
					onclick={() => {
						toast.promise(
							() =>
								storageClient.validateSMBUser({
									realm: realm,
									username: joinSource.username,
									password: joinSource.password,
									searchUsername: request,
									tls: true
								}),
							{
								loading: `Validating ${request}...`,
								success: () => {
									validUsers = [...validUsers, request];
									return `Validated ${request} successfully`;
								},
								error: (error) => {
									let message = `Fail to validate ${request}`;
									toast.error(message, {
										description: (error as ConnectError).message.toString(),
										duration: Number.POSITIVE_INFINITY
									});
									return message;
								}
							}
						);
					}}
				>
					{m.add()}
				</Button>
				<Modal.Action
					disabled={isNull}
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
