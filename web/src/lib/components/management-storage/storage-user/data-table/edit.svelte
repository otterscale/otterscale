<script lang="ts" module>
	import type { UpdateUserRequest, User } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type Writable } from 'svelte/store';
	import { USER_SUSPENDED_HELP_TEXT } from './helper';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		user,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		user: User;
		data: Writable<User[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		userId: user.id,
		userName: user.name,
		suspended: true
	} as UpdateUserRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Edit User
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>ID</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userId} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-placement">Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userName} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-placement">Suspended</Form.Label>
					<SingleInput.Boolean required bind:value={request.suspended} />
				</Form.Field>
				<Form.Help>
					{USER_SUSPENDED_HELP_TEXT}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.updateUser(request)
							.then((r) => {
								toast.success(`Update ${r.name}`);
								storageClient
									.listUsers({ scopeUuid: selectedScope, facilityName: selectedFacility })
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
					}}
				>
					Update
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
