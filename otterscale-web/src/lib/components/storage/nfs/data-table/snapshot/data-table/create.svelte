<script lang="ts" module>
	import type { CreateSubvolumeSnapshotRequest, Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroup,
		subvolume,
		subvolumes: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
		subvolume: Subvolume;
		subvolumes: Writable<Subvolume[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: selectedFacility,
		volumeName: selectedVolume,
		groupName: selectedSubvolumeGroup,
		subvolumeName: subvolume.name
	} as CreateSubvolumeSnapshotRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let invalid = $state(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header>Create Snapshot</AlertDialog.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General id="name" required type="text" bind:value={request.snapshotName} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						stateController.close();
						storageClient
							.createSubvolumeSnapshot(request)
							.then((r) => {
								toast.success(`Create ${r.name}`);
								storageClient
									.listSubvolumes({
										scopeUuid: selectedScopeUuid,
										facilityName: selectedFacility,
										volumeName: selectedVolume,
										groupName: selectedSubvolumeGroup
									})
									.then((r) => {
										data.set(r.subvolumes);
									});
							})
							.catch((e) => {
								toast.error(`Fail to create snapshot: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
