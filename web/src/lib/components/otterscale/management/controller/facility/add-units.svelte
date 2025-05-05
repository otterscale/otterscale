<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import {
		Nexus,
		type AddFacilityUnitsRequest,
		type Facility,
		type Machine,
		type Machine_Placement
	} from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { writable } from 'svelte/store';

	let {
		scopeUuid,
		facilityByCategory = $bindable()
	}: {
		scopeUuid: string;
		facilityByCategory: Facility;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	const DEFAULT_PLACEMENTS = [] as Machine_Placement[];
	const DEFAULT_REQUEST = {
		scopeUuid: scopeUuid,
		name: facilityByCategory.name,
		placements: DEFAULT_PLACEMENTS
	} as AddFacilityUnitsRequest;

	let addFacilityUnitsPlacements = $state(DEFAULT_PLACEMENTS);
	let addFacilityUnitsRequest = $state(DEFAULT_REQUEST);

	function reset() {
		addFacilityUnitsPlacements = DEFAULT_PLACEMENTS;
		addFacilityUnitsRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:plus" />
		Add
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add Units for {facilityByCategory.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid max-h-[77vh] w-full gap-4 overflow-y-auto">
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<Label class="text-sm font-medium">Machines</Label>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger class="w-full">
								<Button variant="outline" class="w-full justify-between">
									<span>Select Machines</span>
									<Icon icon="lucide:chevron-down" />
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end" class="max-h-[300px] overflow-auto">
								{#each $machinesStore as machine}
									<DropdownMenu.Sub>
										<DropdownMenu.SubTrigger>{machine.fqdn}</DropdownMenu.SubTrigger>
										<DropdownMenu.SubContent>
											{#each ['lxd', 'kvm', 'machine'] as type}
												<DropdownMenu.Item
													onclick={() => {
														addFacilityUnitsPlacements = [
															...addFacilityUnitsPlacements,
															{
																machineId: machine.id,
																type: {
																	case: type,
																	value: true
																}
															} as Machine_Placement
														];
													}}
												>
													{type}
												</DropdownMenu.Item>
											{/each}
										</DropdownMenu.SubContent>
									</DropdownMenu.Sub>
								{/each}
							</DropdownMenu.Content>
						</DropdownMenu.Root>
						<div class="flex flex-wrap gap-2">
							{#each addFacilityUnitsPlacements as placement}
								<Badge
									variant="outline"
									onclick={() => {
										addFacilityUnitsPlacements = addFacilityUnitsPlacements.filter(
											(p) => p.machineId !== placement.machineId
										);
									}}
								>
									<div class="flex w-fit items-center gap-1">
										{placement.machineId}: {placement.type.case}
										<Icon icon="ph:x" />
									</div>
								</Badge>
							{/each}
						</div>

						<Label class="text-sm font-medium">Number</Label>
						<Input bind:value={addFacilityUnitsRequest.number} type="number" />
					</fieldset>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					console.log(addFacilityUnitsRequest);
					client
						.addFacilityUnits(addFacilityUnitsRequest)
						.then((r) => {
							toast.info(`Add units to ${addFacilityUnitsRequest.name}`);
							client
								.getFacility({ scopeUuid: scopeUuid, name: facilityByCategory.name })
								.then((r) => {
									facilityByCategory = r;
								});
						})
						.catch((e) => {
							toast.error(`Fail to add units to ${addFacilityUnitsRequest.name}: ${e.toString()}`);
						});
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
