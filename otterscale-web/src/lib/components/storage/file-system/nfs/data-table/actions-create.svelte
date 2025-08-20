<script lang="ts" module>
	import type { CreateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
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
	import { type NFSStore } from '../utils.svelte.js';
	import { SUBVOLUME_QUOTA_HELP_TEXT } from './helper.js';
</script>

<script lang="ts">
	const nfsStore: NFSStore = getContext('nfsStore');
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let open = $state(false);
	function close() {
		open = false;
	}
	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);

	const defaults = {
		scopeUuid: get(nfsStore.selectedScopeUuid),
		facilityName: get(nfsStore.selectedFacilityName),
		volumeName: get(nfsStore.selectedVolumeName),
		groupName: get(nfsStore.selectedSubvolumeGroupName),
		export: true
	} as CreateSubvolumeRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}
</script>

<Modal.Root bind:open>
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
					<SingleInput.General id="name" required type="text" bind:value={request.subvolumeName} />
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
						bind:value={request.quotaBytes}
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
					reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createSubvolume(request), {
							loading: `Creating ${request.subvolumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.subvolumeName}`;
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
