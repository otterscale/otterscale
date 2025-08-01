<script lang="ts" module>
	import type {
		Image,
		Image_Snapshot,
		RollbackImageSnapshotRequest
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
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as RollbackImageSnapshotRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:lock-open" />
		Rollback
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Rollback RADOS Block Device Snapshot</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Label>Name</Form.Label>
				<Form.Help>Please check the snapshot name exactly to confirm this action.</Form.Help>
				<Form.Field>
					<SingleInput.General type="text" disabled bind:value={request.snapshotName} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						toast.info(`Rollbacking ${request.snapshotName}...`);
						storageClient
							.rollbackImageSnapshot(request)
							.then((r) => {
								toast.success(`Rollback ${request.snapshotName}`);
								storageClient
									.listImages({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.images);
									});
							})
							.catch((e) => {
								toast.error(`Fail to rollback snapshot: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Rollback
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
