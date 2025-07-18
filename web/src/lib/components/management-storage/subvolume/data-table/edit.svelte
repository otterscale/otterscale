<script lang="ts" module>
	import type { Subvolume, UpdateSubvolumeRequest } from '$gen/api/storage/v1/storage_pb';
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
	import { SUBVOLUME_QUOTA_HELP_TEXT } from './helper';
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
		groupName: selectedSubvolumeGroup,
		subvolumeName: subvolume.name,
		quotaBytes: subvolume.quotaBytes
	} as UpdateSubvolumeRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Edit Subvolume</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 0), label: 'B' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 1), label: 'KB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 2), label: 'MB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 5), label: 'PB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Help>
					{SUBVOLUME_QUOTA_HELP_TEXT}
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
							.updateSubvolume(request)
							.then((r) => {
								toast.success(`Update ${r.name}`);
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
								toast.error(`Fail to update subvolume: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
					}}
				>
					Update
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
