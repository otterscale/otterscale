<script lang="ts" module>
	import type {
		DeleteSubvolumeSnapshotRequest,
		Subvolume,
		Subvolume_Snapshot
	} from '$lib/api/storage/v1/storage_pb';
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
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		snapshot,
		data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		snapshot: Subvolume_Snapshot;
		data: Writable<Subvolume[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: selectedFacility,
		volumeName: selectedVolume,
		groupName: selectedSubvolumeGroup,
		subvolumeName: subvolume.name
	} as DeleteSubvolumeSnapshotRequest;

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
	<AlertDialog.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Delete Snapshot</AlertDialog.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.DeletionConfirm
						id="deletion"
						required
						target={snapshot.name}
						bind:value={request.snapshotName}
					/>
				</Form.Field>
				<Form.Help>
					Please type the snapshot name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						stateController.close();
						storageClient
							.deleteSubvolumeSnapshot(request)
							.then((r) => {
								toast.success(`Delete ${request.snapshotName}`);
								storageClient
									.listSubvolumes({
										scopeUuid: selectedScopeUuid,
										facilityName: selectedFacility,
										volumeName: selectedVolume,
										groupName: selectedSubvolumeGroup
									})
									.then((r) => {
										data.set(r.subvolumes);
									});
							})
							.catch((e) => {
								toast.error(`Fail to delete snapshot: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
