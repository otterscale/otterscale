<script lang="ts" module>
	import type { DeleteSubvolumeGroupRequest, SubvolumeGroup } from '$gen/api/storage/v1/storage_pb';
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
		subvolumeGroup,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		subvolumeGroup: SubvolumeGroup;
		data: Writable<SubvolumeGroup[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		volumeName: selectedVolume
	} as DeleteSubvolumeGroupRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Delete Subvolume Group
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.DeletionConfirm
						required
						target={subvolumeGroup.name}
						bind:value={request.groupName}
					/>
				</Form.Field>
				<Form.Help>
					Please type the file subvolume group exactly to confirm deletion. This action cannot be
					undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.deleteSubvolumeGroup(request)
							.then((r) => {
								toast.success(`Delete ${request.groupName}`);
								storageClient
									.listSubvolumeGroups({
										scopeUuid: selectedScope,
										facilityName: selectedFacility,
										volumeName: selectedVolume
									})
									.then((r) => {
										data.set(r.subvolumeGroups);
									});
							})
							.catch((e) => {
								toast.error(`Fail to delete subvolume group: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
					}}
				>
					Delete
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
