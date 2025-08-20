<script lang="ts">
	import type {
		Image,
		Image_Snapshot,
		RollbackImageSnapshotRequest
	} from '$lib/api/storage/v1/storage_pb';
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

	let {
		snapshot
	}: {
		snapshot: Image_Snapshot;
	} = $props();

	const image: Image = getContext('image');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		imageName: image.name,
		poolName: image.poolName
	} as RollbackImageSnapshotRequest;
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
					<SingleInput.Confirm
						required
						target={snapshot.name}
						bind:value={request.snapshotName}
						bind:invalid
					/>
				</Form.Field>
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
						toast.promise(() => storageClient.rollbackImageSnapshot(request), {
							loading: `Rolling back ${request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Rollback ${request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to rollback ${request.snapshotName}`;
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
					Rollback
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
