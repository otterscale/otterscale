<script lang="ts" module>
	import type { SubvolumeGroup, UpdateSubvolumeGroupRequest } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
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
		volumeName: selectedVolume,
		groupName: subvolumeGroup.name,
		quotaBytes: subvolumeGroup.quotaBytes
	} as UpdateSubvolumeGroupRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const poolOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	let isPoolsLoading = $state(true);
	async function fetchPools() {
		try {
			const response = await storageClient.listPools({
				scopeUuid: selectedScope,
				facilityName: selectedFacility
			});

			poolOptions.set(
				response.pools.map((pool) => ({
					value: pool.name,
					label: pool.name,
					icon: 'ph:cube'
				}))
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isPoolsLoading = false;
		}
	}

	onMount(async () => {
		try {
			await fetchPools();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Edit Subvolume Group
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.groupName} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<SingleInput.General bind:value={request.quotaBytes} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.updateSubvolumeGroup(request)
							.then((r) => {
								toast.success(`Update ${r.name}`);
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
								toast.error(`Fail to update subvolume group: ${e.toString()}`);
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
