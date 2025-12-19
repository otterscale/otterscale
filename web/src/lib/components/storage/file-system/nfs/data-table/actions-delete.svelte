<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DeleteSubvolumeRequest, Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let {
		subvolume,
		scope,
		volume,
		group,
		reloadManager,
		closeActions
	}: {
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');

	let open = $state(false);
	function close() {
		open = false;
	}
	let invalid = $state(false);

	let request = $state({} as DeleteSubvolumeRequest);
	function init() {
		request = {
			scope: scope,
			volumeName: volume,
			groupName: group
		} as DeleteSubvolumeRequest;
	}

	const storageClient = createClient(StorageService, transport);
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
		<Modal.Header>
			{m.delete_nfs()}
		</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						required
						target={subvolume.name}
						bind:value={request.subvolumeName}
						bind:invalid
					/>
				</Form.Field>
				<Form.Help>
					{m.deletion_warning({ identifier: m.nfs_name() })}
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
						toast.promise(() => storageClient.deleteSubvolume(request), {
							loading: `Deleting ${request.subvolumeName}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.subvolumeName}`;
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
