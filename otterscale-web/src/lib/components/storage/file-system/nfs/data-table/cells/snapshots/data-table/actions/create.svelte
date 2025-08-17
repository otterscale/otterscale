<script lang="ts" module>
	import type { CreateSubvolumeSnapshotRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Modal from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import type { nfsSnapshotStore } from '../../utils.svelte';
</script>

<script lang="ts">
	const nfsSnapshotStore: nfsSnapshotStore = getContext('nfsSnapshotStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid = $state(false);

	const storageClient = createClient(StorageService, transport);
	const requestManager = new RequestManager<CreateSubvolumeSnapshotRequest>({
		scopeUuid: get(nfsSnapshotStore.selectedScopeUuid),
		facilityName: get(nfsSnapshotStore.selectedFacilityName),
		volumeName: get(nfsSnapshotStore.selectedVolumeName),
		groupName: get(nfsSnapshotStore.selectedSubvolumeGroupName),
		subvolumeName: get(nfsSnapshotStore.selectedSubvolumeName)
	} as CreateSubvolumeSnapshotRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<Modal.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</Modal.Trigger>
	</div>
	<Modal.Content>
		<Modal.Header>Create Snapshot</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
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
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createSubvolumeSnapshot(requestManager.request), {
							loading: `Creating ${requestManager.request.snapshotName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${requestManager.request.snapshotName}`;
							},
							error: (error) => {
								let message = `Fail to create ${requestManager.request.snapshotName}`;
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
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
