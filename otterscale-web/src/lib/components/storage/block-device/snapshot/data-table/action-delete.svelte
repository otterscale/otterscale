<script lang="ts" module>
	import type {
		DeleteImageSnapshotRequest,
		Image,
		Image_Snapshot
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

	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		imageName: image.name,
		poolName: image.poolName
	} as DeleteImageSnapshotRequest;
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
		<Modal.Header>Delete RADOS Block Device Snapshot</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						required
						target={snapshot.name}
						bind:value={request.snapshotName}
						bind:invalid
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
					reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.deleteImageSnapshot(request), {
							loading: `Deleting ${request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Delete ${request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.snapshotName}`;
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
