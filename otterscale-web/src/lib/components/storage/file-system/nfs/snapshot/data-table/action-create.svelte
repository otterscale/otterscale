<script lang="ts">
	import type { CreateSubvolumeSnapshotRequest, Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import type { NFSStore } from '../../utils.svelte';

	const nfsStore: NFSStore = getContext('nfsStore');
	const subvolume: Subvolume = getContext('subvolume');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid = $state(false);

	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scopeUuid: get(nfsStore.selectedScopeUuid),
		facilityName: get(nfsStore.selectedFacilityName),
		volumeName: get(nfsStore.selectedVolumeName),
		groupName: get(nfsStore.selectedSubvolumeGroupName),
		subvolumeName: subvolume.name
	} as CreateSubvolumeSnapshotRequest;
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
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Snapshot</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General id="name" required type="text" bind:value={request.snapshotName} />
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
						toast.promise(() => storageClient.createSubvolumeSnapshot(request), {
							loading: `Creating ${request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.snapshotName}`;
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
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
