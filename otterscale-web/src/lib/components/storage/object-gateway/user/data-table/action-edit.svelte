<script lang="ts" module>
	import type { UpdateUserRequest, User } from '$lib/api/storage/v1/storage_pb';
	import { m } from '$lib/paraglide/messages.js';
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
	import { user_suspended_descriptor } from './helper';
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

	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		userId: user.id,
		userName: user.name,
		suspended: true
	} as UpdateUserRequest;
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
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_user()}</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.id()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userId} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userName} />
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean
						format="checkbox"
						descriptor={user_suspended_descriptor}
						bind:value={request.suspended}
					/>
					<Form.Help>
						{m.user_suspended_direction()}
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.updateUser(request), {
							loading: `Updating ${request.userId}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${request.userId}`;
							},
							error: (error) => {
								let message = `Fail to update ${request.userId}`;
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
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
