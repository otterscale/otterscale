<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

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
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		snapshot,
		image,
		scope,
		reloadManager,
		closeActions
	}: {
		snapshot: Image_Snapshot;
		image: Image;
		scope: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');

	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);

	let request = $state({} as RollbackImageSnapshotRequest);
	function init() {
		request = {
			scope: scope,
			imageName: image.name,
			poolName: image.poolName
		} as RollbackImageSnapshotRequest;
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
		<Icon icon="ph:arrow-counter-clockwise" />
		{m.rollback()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.rollback_snapshot()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Label>{m.name()}</Form.Label>
				<Form.Help>{m.deletion_warning({ identifier: m.snapshot_name() })}</Form.Help>
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
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.rollbackImageSnapshot(request), {
							loading: `Rolling back ${request.snapshotName}...`,
							success: () => {
								reloadManager.force();
								return `Rollback ${request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to rollback ${request.snapshotName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
									closeButton: true
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
