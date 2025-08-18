<script lang="ts" module>
	import type { CreateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
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
	const nfsStore: NFSStore = getContext('nfsStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const stateController = new StateController(false);
	let invalid = $state(false);

	const requestManager = new RequestManager<CreateSubvolumeRequest>({
		scopeUuid: get(nfsStore.selectedScopeUuid),
		facilityName: get(nfsStore.selectedFacilityName),
		volumeName: get(nfsStore.selectedVolumeName),
		groupName: get(nfsStore.selectedSubvolumeGroupName),
		export: true
	} as CreateSubvolumeRequest);
	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header class="flex items-center justify-center text-xl font-bold">
			Create NFS
		</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
						bind:value={requestManager.request.subvolumeName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Group</Form.Label>
					<SingleInput.General
						required
						disabled
						type="text"
						value={get(nfsStore.selectedSubvolumeGroupName) === ''
							? 'default'
							: get(nfsStore.selectedSubvolumeGroupName)}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Size</Form.Label>
					<Form.Help>
						{SUBVOLUME_QUOTA_HELP_TEXT}
					</Form.Help>
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
						toast.promise(() => storageClient.createSubvolume(requestManager.request), {
							loading: `Creating ${requestManager.request.subvolumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${requestManager.request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to create ${requestManager.request.subvolumeName}`;
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
