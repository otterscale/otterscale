<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

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
		subvolume,
		scope,
		volume,
		reloadManager,
		closeActions
	}: {
		subvolume: Subvolume;
		scope: string;
		volume: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();
	const transport: Transport = getContext('transport');

	let open = $state(false);
	function close() {
		open = false;
	}
	let invalid = $state(false);

	let request = $state({} as RevokeSubvolumeExportAccessRequest);
	function init() {
		request = {
			scope: scope,
			volumeName: volume,
			subvolumeName: subvolume.name
		} as RevokeSubvolumeExportAccessRequest;
	}

	const storageClient = createClient(StorageService, transport);
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:shield-slash" />
		{m.revoke()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.revoke_export_access()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.client_ip()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.clientIp} bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
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
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
