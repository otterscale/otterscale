<script lang="ts" module>
	import type { SubvolumeGroup, UpdateSubvolumeGroupRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import { type GroupStore } from '../../utils.svelte.js';
</script>

<script lang="ts">
	let {
		subvolumeGroup
	}: {
		subvolumeGroup: SubvolumeGroup;
	} = $props();

	const groupStore: GroupStore = getContext('groupStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const requestManager = new RequestManager<UpdateSubvolumeGroupRequest>({
		scopeUuid: get(groupStore.selectedScopeUuid),
		facilityName: get(groupStore.selectedFacilityName),
		volumeName: get(groupStore.selectedVolumeName),
		groupName: subvolumeGroup.name,
		quotaBytes: subvolumeGroup.quotaBytes
	} as UpdateSubvolumeGroupRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Subvolume Group</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
						bind:value={requestManager.request.groupName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Size</Form.Label>
					<SingleInput.Measurement
						bind:value={requestManager.request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
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
						toast.promise(() => storageClient.updateSubvolumeGroup(requestManager.request), {
							loading: `Updating ${requestManager.request.volumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${requestManager.request.volumeName}`;
							},
							error: (error) => {
								let message = `Fail to update ${requestManager.request.volumeName}`;
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
					Update
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
