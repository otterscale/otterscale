<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type {
		DeleteSubvolumeSnapshotRequest,
		Subvolume,
		Subvolume_Snapshot
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		snapshot,
		subvolume,
		scope,
		volume,
		group,
		reloadManager,
		closeActions
	}: {
		snapshot: Subvolume_Snapshot;
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let invalid = $state(false);

	let request = $state({} as DeleteSubvolumeSnapshotRequest);
	function init() {
		request = {
			scope: scope,
			volumeName: volume,
			groupName: group,
			subvolumeName: subvolume.name
		} as DeleteSubvolumeSnapshotRequest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_snapshot()}</Modal.Header>
		<Form.Root>
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
					{m.deletion_warning({ identifier: m.snapshot_name() })}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.deleteSubvolumeSnapshot(request), {
							loading: `Deleting ${request.snapshotName}...`,
							success: () => {
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
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
