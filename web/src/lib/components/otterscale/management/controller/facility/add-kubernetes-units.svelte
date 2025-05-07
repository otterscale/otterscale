<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Nexus,
		type AddKubernetesUnitsRequest,
		type Facility,
		type Machine
	} from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';

	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { writable } from 'svelte/store';

	let {
		scopeUuid,
		kubernetes
	}: {
		scopeUuid: string;
		kubernetes: Facility;
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

	const DEFAULT_MACHINES = [] as Machine[];
	const DEFAULT_REQUEST = {
		scopeUuid: scopeUuid,
		facilityName: kubernetes.name,
		machineIds: [] as string[],
		force: false
	} as AddKubernetesUnitsRequest;

	let addKubernetesUnitsRequest = $state(DEFAULT_REQUEST);
	let selectedMachines = $state(DEFAULT_MACHINES);

	function reset() {
		addKubernetesUnitsRequest = DEFAULT_REQUEST;
		selectedMachines = DEFAULT_MACHINES;
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
			23;
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Button class="ml-auto" variant="ghost">
			<Icon icon="ph:plus" />
			Bundle Units
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add Units for {kubernetes.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid max-h-[77vh] w-full gap-4 overflow-y-auto">
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<Label class="text-sm font-medium">Number</Label>
						<Input bind:value={addKubernetesUnitsRequest.number} type="number" />

						<Label class="text-sm font-medium">Machines</Label>
						{#if selectedMachines.length > 0}
							<span class="flex flex-wrap items-center gap-1">
								{#each selectedMachines as machine}
									<Badge variant="secondary">{machine.fqdn}</Badge>
								{/each}
							</span>
						{/if}
						<Select.Root type="multiple" bind:value={addKubernetesUnitsRequest.machineIds}>
							<Select.Trigger>Select</Select.Trigger>
							<Select.Content class="max-h-[230px] overflow-y-auto" sideOffset={7}>
								{#each $machinesStore as machine}
									<Select.Item
										value={machine.id}
										onclick={() => {
											if (!selectedMachines.includes(machine)) {
												selectedMachines = [...selectedMachines, machine];
											} else {
												selectedMachines = selectedMachines.filter((m) => m.id !== machine.id);
											}
										}}
									>
										{machine.fqdn}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>

						<div class="flex items-center justify-between">
							<Label class="text-sm font-medium">Force</Label>
							<Switch bind:checked={addKubernetesUnitsRequest.force} />
						</div>
						<div class="flex flex-col gap-2 rounded-lg bg-muted p-2">
							<p class="text-sm">Note that more than 3 units will slow down High Availability</p>
							<p class="text-xs font-light text-muted-foreground">
								To add additional worker nodes, scroll down to find the worker section and click the
								"Add Unit" button.
							</p>
						</div>
					</fieldset>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					toast.promise(() => client.addKubernetesUnits(addKubernetesUnitsRequest), {
						loading: 'Loading...',
						success: (r) => {
							return `Add units for ${addKubernetesUnitsRequest.facilityName} success`;
						},
						error: (e) => {
							let msg = `Fail to add units for ${addKubernetesUnitsRequest.facilityName}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
						}
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
