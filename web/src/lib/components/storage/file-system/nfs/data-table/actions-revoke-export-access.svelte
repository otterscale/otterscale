<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { type NFSStore } from '../utils.svelte.js';

	import type {
		RevokeSubvolumeExportAccessRequest,
		Subvolume
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
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

	let open = $state(false);
	function close() {
		open = false;
	}
	let invalid = $state(false);

	const defaults = {
		scope: get(nfsStore.selectedScope),
		facility: get(nfsStore.selectedFacility),
		volumeName: get(nfsStore.selectedVolumeName),
		subvolumeName: subvolume.name
	} as RevokeSubvolumeExportAccessRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:shield-slash" />
		{m.revoke()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.revoke_export_access()}</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.client_ip()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.clientIp} />
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
						toast.promise(() => storageClient.revokeSubvolumeExportAccess(request), {
							loading: `Revoking ${request.subvolumeName}...`,
							success: () => {
								reloadManager.force();
								return `Revoke ${request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to revoke ${request.subvolumeName}`;
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
