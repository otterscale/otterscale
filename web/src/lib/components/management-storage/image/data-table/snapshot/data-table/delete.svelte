<script lang="ts" module>
	import type {
		DeleteImageSnapshotRequest,
		Image,
		Image_Snapshot
	} from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		image,
		snapshot,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		image: Image;
		snapshot: Image_Snapshot;
		data: Writable<Image[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		imageName: image.name,
		poolName: image.poolName
	} as DeleteImageSnapshotRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

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
		<AlertDialog.Header>Delete RADOS Block Device Snapshot</AlertDialog.Header>
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
						toast.info(`Deleting ${request.snapshotName}...`);
						storageClient
							.deleteImageSnapshot(request)
							.then((r) => {
								toast.success(`Delete ${request.snapshotName}`);
								storageClient
									.listImages({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.images);
									});
							})
							.catch((e) => {
								toast.error(`Fail to delete snapshot: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
