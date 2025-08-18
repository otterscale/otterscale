<script lang="ts" module>
	import type { UpdateUserRequest, User } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
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
	let {
		user
	}: {
		user: User;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const requestManager = new RequestManager<UpdateUserRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		userId: user.id,
		userName: user.name,
		suspended: true
	} as UpdateUserRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit User</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>ID</Form.Label>
					<SingleInput.General
						id="id"
						required
						type="text"
						bind:value={requestManager.request.userId}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
						bind:value={requestManager.request.userName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Suspended</Form.Label>
					<Form.Help>
						{USER_SUSPENDED_HELP_TEXT}
					</Form.Help>
					<SingleInput.Boolean
						format="checkbox"
						descriptor={user_suspended_descriptor}
						bind:value={requestManager.request.suspended}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.updateUser(requestManager.request), {
							loading: `Updating ${requestManager.request.userId}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${requestManager.request.userId}`;
							},
							error: (error) => {
								let message = `Fail to update ${requestManager.request.userId}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						requestManager.reset();
						stateController.close();
					}}
				>
					Update
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
