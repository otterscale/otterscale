<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { Subvolume, UpdateSubvolumeRequest } from '$lib/api/storage/v1/storage_pb';
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
		reloadManager
	}: {
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	let open = $state(false);
	function close() {
		open = false;
	}

	const defaults = {
		scope: scope,
		volumeName: volume,
		groupName: group,
		subvolumeName: subvolume.name,
		quotaBytes: subvolume.quotaBytes
	} as UpdateSubvolumeRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_nfs()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>{m.quota_size()}</Form.Legend>

				<Form.Field>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Help>
					{m.nfs_quota_size_direction()}
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
					onclick={() => {
						toast.promise(() => storageClient.updateSubvolume(request), {
							loading: `Updating ${request.subvolumeName}...`,
							success: () => {
								reloadManager.force();
								return `Update ${request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to update ${request.subvolumeName}`;
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
