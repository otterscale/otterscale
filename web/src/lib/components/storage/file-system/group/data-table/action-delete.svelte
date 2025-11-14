<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { DeleteSubvolumeGroupRequest, SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';

	import { type GroupStore } from '../utils.svelte';
</script>

<script lang="ts">
	let {
		subvolumeGroup
	}: {
		subvolumeGroup: SubvolumeGroup;
	} = $props();

	const groupStore: GroupStore = getContext('groupStore');
	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const defaults = {
		scope: get(groupStore.selectedScope),
		facility: get(groupStore.selectedFacility),
		volumeName: get(groupStore.selectedVolumeName)
	} as DeleteSubvolumeGroupRequest;
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
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_group()}</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						id="deletion"
						required
						target={subvolumeGroup.name}
						bind:value={request.groupName}
					/>
				</Form.Field>
				<Form.Help>
					{m.deletion_warning({ identifier: m.group_name() })}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.deleteSubvolumeGroup(request), {
							loading: `Deleting ${request.volumeName}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${request.volumeName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.volumeName}`;
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
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
