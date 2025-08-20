<script lang="ts" module>
	import type { CreateUserRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { USER_SUSPENDED_HELP_TEXT, user_suspended_descriptor } from './helper';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		suspended: true
	} as CreateUserRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create User</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>ID</Form.Label>
					<SingleInput.General id="id" required type="text" bind:value={request.userId} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General id="name" required type="text" bind:value={request.userName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Suspended</Form.Label>
					<Form.Help>
						{USER_SUSPENDED_HELP_TEXT}
					</Form.Help>
					<SingleInput.Boolean
						format="checkbox"
						descriptor={user_suspended_descriptor}
						bind:value={request.suspended}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createUser(request), {
							loading: `Creating ${request.userName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${request.userName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.userName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						reset();
						close();
					}}
				>
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
