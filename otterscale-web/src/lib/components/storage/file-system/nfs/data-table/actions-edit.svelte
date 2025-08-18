<script lang="ts" module>
	import type { Subvolume, UpdateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import { type NFSStore } from '../utils.svelte.js';
	import { SUBVOLUME_QUOTA_HELP_TEXT } from './helper.js';
</script>

<script lang="ts">
	let {
		subvolume
	}: {
		subvolume: Subvolume;
	} = $props();

	const nfsStore: NFSStore = getContext('nfsStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const stateController = new StateController(false);

	const requestManager = new RequestManager<UpdateSubvolumeRequest>({
		scopeUuid: get(nfsStore.selectedScopeUuid),
		facilityName: get(nfsStore.selectedFacilityName),
		volumeName: get(nfsStore.selectedVolumeName),
		groupName: get(nfsStore.selectedSubvolumeGroupName),
		subvolumeName: subvolume.name,
		quotaBytes: subvolume.quotaBytes
	} as UpdateSubvolumeRequest);
	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Subvolume</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>Quota</Form.Legend>

				<Form.Field>
					<SingleInput.Measurement
						bind:value={requestManager.request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Help>
					{SUBVOLUME_QUOTA_HELP_TEXT}
				</Form.Help>
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
						toast.promise(() => storageClient.updateSubvolume(requestManager.request), {
							loading: `Updating ${requestManager.request.subvolumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${requestManager.request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to update ${requestManager.request.subvolumeName}`;
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
