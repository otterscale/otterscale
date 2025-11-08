<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';

	import { type NFSStore } from '../utils.svelte.js';
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
		scope: get(nfsStore.selectedScope),
		facility: get(nfsStore.selectedFacility),
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
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header class="flex items-center justify-center text-xl font-bold">
			{m.create_nfs()}
		</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.subvolumeName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.group()}</Form.Label>
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
					<Form.Label>{m.quota_size()}</Form.Label>
					<Form.Help>
						{m.nfs_quota_size_direction()}
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
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createSubvolume(request), {
							loading: `Creating ${request.subvolumeName}...`,
							success: () => {
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
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
