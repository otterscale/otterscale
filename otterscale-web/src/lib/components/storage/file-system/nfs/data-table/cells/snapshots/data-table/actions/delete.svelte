<script lang="ts" module>
	import type {
		DeleteSubvolumeSnapshotRequest,
		Subvolume_Snapshot
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Modal from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import type { nfsSnapshotStore } from '../../utils.svelte';
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Subvolume_Snapshot;
	} = $props();

	const nfsSnapshotStore: nfsSnapshotStore = getContext('nfsSnapshotStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid = $state(false);

	const storageClient = createClient(StorageService, transport);
	const requestManager = new RequestManager<DeleteSubvolumeSnapshotRequest>({
		scopeUuid: get(nfsSnapshotStore.selectedScopeUuid),
		facilityName: get(nfsSnapshotStore.selectedFacilityName),
		volumeName: get(nfsSnapshotStore.selectedVolumeName),
		groupName: get(nfsSnapshotStore.selectedSubvolumeGroupName),
		subvolumeName: get(nfsSnapshotStore.selectedSubvolumeName)
	} as DeleteSubvolumeSnapshotRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Snapshot</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						id="deletion"
						required
						target={snapshot.name}
						bind:value={requestManager.request.snapshotName}
					/>
				</Form.Field>
				<Form.Help>
					Please type the snapshot name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
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
						toast.promise(() => storageClient.deleteSubvolumeSnapshot(requestManager.request), {
							loading: `Deleting ${requestManager.request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Delete ${requestManager.request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${requestManager.request.snapshotName}`;
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
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
