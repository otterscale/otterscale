<script lang="ts" module>
	import type { UpdateUserRequest, User } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type Writable } from 'svelte/store';
	import { USER_SUSPENDED_HELP_TEXT, user_suspended_descriptor } from './helper';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		user,
		data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		user: User;
		data: Writable<User[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: selectedFacility,
		userId: user.id,
		userName: user.name,
		suspended: true
	} as UpdateUserRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let invalid = $state(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Edit User</AlertDialog.Header>
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
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						toast.info(`Updating ${request.userId}...`);
						storageClient
							.updateUser(request)
							.then((r) => {
								toast.success(`Update ${r.name}`);
								storageClient
									.listUsers({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.users);
									});
							})
							.catch((e) => {
								toast.error(`Fail to update user: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Update
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
