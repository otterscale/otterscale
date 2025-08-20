<script lang="ts">
	import type { DeleteSubvolumeGroupRequest, SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import { type GroupStore } from '../utils.svelte';

	let {
		subvolumeGroup
	}: {
		subvolumeGroup: SubvolumeGroup;
	} = $props();

	const groupStore: GroupStore = getContext('groupStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const defaults = {
		scopeUuid: get(groupStore.selectedScopeUuid),
		facilityName: get(groupStore.selectedFacilityName),
		volumeName: get(groupStore.selectedVolumeName)
	} as DeleteSubvolumeGroupRequest;
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
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Subvolume Group</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						id="deletion"
						required
						target={subvolumeGroup.name}
						bind:value={request.groupName}
					/>
				</Form.Field>
				<Form.Help>
					Please type the file subvolume group exactly to confirm deletion. This action cannot be
					undone.
				</Form.Help>
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
						toast.promise(() => storageClient.deleteSubvolumeGroup(request), {
							loading: `Deleting ${request.volumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Delete ${request.volumeName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.volumeName}`;
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
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
