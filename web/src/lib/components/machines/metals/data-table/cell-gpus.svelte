<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import { page } from '$app/state';
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import {
		type GPURelation_GPU,
		type GPURelation_Machine,
		type GPURelation_Pod,
		OrchestratorService
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Table from '$lib/components/custom/table';
	import { Complex as Simple } from '$lib/components/flow/index';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const machineScope = $derived(
		machine.workloadAnnotations['juju-machine-id']?.split('-machine-')[0]
	);

	const transport: Transport = getContext('transport');
	const client = createClient(OrchestratorService, transport);

	const position = { x: 0, y: 0 };

	const nodes: Writable<Node[]> = writable([]);
	const edges: Writable<Edge[]> = writable([]);
	let isLoading = $state(false);
	let hasLoadedData = $state(false);

	async function loadGPURelations() {
		if (hasLoadedData) return;

		isLoading = true;
		try {
			const response = await client.listGPURelationsByMachine({
				scope: machineScope,
				machineId: machine.id
			});

			// Extract GPUs and pods once to avoid redundant computation
			const gpus = response.gpuRelations
				.filter((gpuRelation) => gpuRelation.entity.case === 'gpu')
				.map((gpuRelation) => gpuRelation.entity.value as GPURelation_GPU);

			const pods = response.gpuRelations
				.filter((gpuRelation) => gpuRelation.entity.case === 'pod')
				.map((gpuRelation) => gpuRelation.entity.value as GPURelation_Pod);

			nodes.set(
				response.gpuRelations.map((gpuRelation) => {
					if (gpuRelation.entity.case === 'machine') {
						return {
							type: 'machine',
							id: `machine${gpuRelation.entity.value.id}`,
							data: {
								machineScope,
								machine: gpuRelation.entity.value,
								gpus: gpus.filter(
									(gpu) => gpu.machineId === (gpuRelation.entity.value as GPURelation_Machine).id
								)
							},
							position
						};
					} else if (gpuRelation.entity.case === 'gpu') {
						return {
							type: 'gpu',
							id: `gpu${gpuRelation.entity.value.id}`,
							data: {
								machineScope,
								gpu: gpuRelation.entity.value,
								devices: pods.flatMap((pod) =>
									pod.devices.filter((device) => {
										return device.gpuId === (gpuRelation.entity.value as GPURelation_GPU).id;
									})
								)
							},
							position
						};
					} else if (gpuRelation.entity.case === 'pod') {
						return {
							type: 'model',
							id: `pod${gpuRelation.entity.value.namespace}${gpuRelation.entity.value.name}`,
							data: gpuRelation.entity.value,
							position
						};
					} else {
						return {} as Node;
					}
				})
			);

			edges.set(
				response.gpuRelations.flatMap((gpuRelation) => {
					if (gpuRelation.entity.case === 'gpu') {
						const gpu = gpuRelation.entity.value as GPURelation_GPU;
						return [
							{
								type: 'edge',
								id: `gpu${gpu.id}machine${gpu.machineId}`,
								source: `gpu${gpu.id}`,
								target: `machine${gpu.machineId}`,
								animated: true,
								selectable: false
							}
						];
					} else if (gpuRelation.entity.case === 'pod') {
						return gpuRelation.entity.value.devices.map((device) => {
							const pod = gpuRelation.entity.value as GPURelation_Pod;
							return {
								type: 'edge',
								id: `pod${pod.namespace}${pod.name}gpu${device.gpuId}`,
								source: `pod${pod.namespace}${pod.name}`,
								target: `gpu${device.gpuId}`,
								animated: true,
								selectable: false
							};
						});
					} else {
						return [];
					}
				})
			);

			hasLoadedData = true;
		} catch (error) {
			console.error('Error loading GPU relations:', error);
		} finally {
			isLoading = false;
		}
	}

	let open = $state(false);

	const isModelEnabled = $derived(page.data['model-enabled']);
</script>

{#if machineScope && isModelEnabled}
	<Sheet.Root bind:open onOpenChange={loadGPURelations}>
		<Sheet.Trigger>
			<span class="flex items-center gap-1">
				{machine.gpuDevices.length}
				<Icon icon="ph:arrow-square-out" />
			</span>
		</Sheet.Trigger>
		<Sheet.Content side="right" class="min-w-[38vw] p-4">
			{#if open}
				<Sheet.Header>
					<Sheet.Title class="text-center text-lg">{m.details()}</Sheet.Title>
				</Sheet.Header>
				{#if isLoading}
					<div class="flex h-full items-center justify-center p-8">
						<Icon icon="ph:circle-notch" class="h-8 w-8 animate-spin" />
					</div>
				{:else}
					<Simple.Flow initialNodes={$nodes} initialEdges={$edges} />
				{/if}
				<div class="mt-auto rounded-lg border shadow">
					<Table.Root>
						<Table.Header>
							<Table.Row
								class="bg-muted [&_th]:bg-muted [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg"
							>
								<Table.Head>{m.product()}</Table.Head>
								<Table.Head>{m.vendor()}</Table.Head>
								<Table.Head>{m.bus()}</Table.Head>
								<Table.Head>{m.pci_address()}</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each machine.gpuDevices as gpuDevice}
								<Table.Row>
									<Table.Cell>
										{gpuDevice.productName !== '' ? gpuDevice.productName : gpuDevice.productId}
									</Table.Cell>
									<Table.Cell>
										{gpuDevice.vendorName !== '' ? gpuDevice.vendorName : gpuDevice.vendorId}
									</Table.Cell>
									<Table.Cell>{gpuDevice.busName}</Table.Cell>
									<Table.Cell>{gpuDevice.pciAddress}</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>
			{/if}
		</Sheet.Content>
	</Sheet.Root>
{:else}
	<HoverCard.Root>
		<HoverCard.Trigger>
			<span class="flex items-center gap-1">
				{machine.gpuDevices.length}
				<Icon icon="ph:info" />
			</span>
		</HoverCard.Trigger>
		<HoverCard.Content class="m-0 h-fit w-fit p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row
						class="bg-muted [&_th]:bg-muted [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg"
					>
						<Table.Head>{m.product()}</Table.Head>
						<Table.Head>{m.vendor()}</Table.Head>
						<Table.Head>{m.bus()}</Table.Head>
						<Table.Head>{m.pci_address()}</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each machine.gpuDevices as gpuDevice}
						<Table.Row>
							<Table.Cell>
								{gpuDevice.productName !== '' ? gpuDevice.productName : gpuDevice.productId}
							</Table.Cell>
							<Table.Cell>
								{gpuDevice.vendorName !== '' ? gpuDevice.vendorName : gpuDevice.vendorId}
							</Table.Cell>
							<Table.Cell>{gpuDevice.busName}</Table.Cell>
							<Table.Cell>{gpuDevice.pciAddress}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/if}
