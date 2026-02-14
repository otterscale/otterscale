<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { CreateUserRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let {
		scope,
		reloadManager
	}: {
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	let request = $state({} as CreateUserRequest);
	function init() {
		request = {
			scope: scope,
			suspended: false
		} as CreateUserRequest;
	}

	let invalidity = $state({} as Booleanified<CreateUserRequest>);
	const invalid = $derived(invalidity.userId || invalidity.userName);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>
			{m.create_user()}
		</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.id()}</Form.Label>
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.userId}
						bind:invalid={invalidity.userId}
						validateRule="lower-alphanum-dash-start-alpha"
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.userName}
						bind:invalid={invalidity.userName}
					/>
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean descriptor={() => m.suspend()} bind:value={request.suspended} />
					<Form.Help>
						{m.user_suspended_direction()}
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createUser(request), {
							loading: `Creating ${request.userName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.userName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.userName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
									closeButton: true
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
