<script lang="ts" module>
	import type {
		Image,
		Image_Snapshot,
		RollbackImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
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
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Image_Snapshot;
	} = $props();

	const image: Image = getContext('image');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	const requestManager = new RequestManager<RollbackImageSnapshotRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		imageName: image.name,
		poolName: image.poolName,
		snapshotName: snapshot.name
	} as RollbackImageSnapshotRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:lock-open" />
		Rollback
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Rollback RADOS Block Device Snapshot</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Label>Name</Form.Label>
				<Form.Help>Please check the snapshot name exactly to confirm this action.</Form.Help>
				<Form.Field>
					<SingleInput.General
						type="text"
						disabled
						bind:value={requestManager.request.snapshotName}
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
					onclick={() => {
						toast.promise(() => storageClient.rollbackImageSnapshot(requestManager.request), {
							loading: `Rolling back ${requestManager.request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Rollback ${requestManager.request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to rollback ${requestManager.request.snapshotName}`;
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
					Rollback
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
