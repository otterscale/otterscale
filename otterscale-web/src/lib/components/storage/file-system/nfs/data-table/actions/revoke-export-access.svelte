<script lang="ts" module>
	import type {
		RevokeSubvolumeExportAccessRequest,
		Subvolume
	} from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { get } from 'svelte/store';
	import { type NFSStore } from '../../utils.svelte.js';
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
	let invalid = $state(false);

	const requestManager = new RequestManager<RevokeSubvolumeExportAccessRequest>({
		scopeUuid: get(nfsStore.selectedScopeUuid),
		facilityName: get(nfsStore.selectedFacilityName),
		volumeName: get(nfsStore.selectedVolumeName),
		subvolumeName: subvolume.name
	} as RevokeSubvolumeExportAccessRequest);
	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:shield-slash" />
		Revoke
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Revoke Export Access</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Client IP</Form.Label>
					<SingleInput.General
						id="client_ip"
						required
						type="text"
						bind:value={requestManager.request.clientIp}
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
						toast.promise(() => storageClient.revokeSubvolumeExportAccess(requestManager.request), {
							loading: `Revoking ${requestManager.request.subvolumeName}...`,
							success: (response) => {
								reloadManager.force();
								return `Revoke ${requestManager.request.subvolumeName}`;
							},
							error: (error) => {
								let message = `Fail to revoke ${requestManager.request.subvolumeName}`;
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
					Revoke
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
