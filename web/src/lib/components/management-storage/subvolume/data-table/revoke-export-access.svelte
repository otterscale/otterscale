<script lang="ts" module>
	import type {
		RevokeSubvolumeExportAccessRequest,
		Subvolume
	} from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		data: Writable<Subvolume[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		volumeName: selectedVolume,
		subvolumeName: subvolume.name
	} as RevokeSubvolumeExportAccessRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let invalid = $state(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:shield-slash" />
		Revoke
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Revoke Export Access</AlertDialog.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Client IP</Form.Label>
					<SingleInput.General id="client_ip" required type="text" bind:value={request.clientIp} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						toast.info(`Revoking ${request.subvolumeName}...`);
						storageClient
							.revokeSubvolumeExportAccess(request)
							.then((r) => {
								toast.success(`Revoke ${request.subvolumeName}`);
								storageClient
									.listSubvolumes({
										scopeUuid: selectedScope,
										facilityName: selectedFacility,
										volumeName: selectedVolume,
										groupName: selectedSubvolumeGroup
									})
									.then((r) => {
										data.set(r.subvolumes);
									});
							})
							.catch((e) => {
								toast.error(`Fail to revoke subvolume: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Revoke
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
