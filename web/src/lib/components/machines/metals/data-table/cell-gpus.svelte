<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';
	import { getContext, onMount } from 'svelte';
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
	import { Complex as ComplexFlow } from '$lib/components/flow/index';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Drawer from '$lib/components/ui/drawer/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card';
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
	let isLoaded = $state(false);

	async function fetch() {
		if (!machineScope) {
			return;
		}

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
						data: { scope: machineScope, pod: gpuRelation.entity.value },
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
	}

	let open = $state(false);

	const isModelEnabled = $derived(page.data['model-enabled']);

	onMount(async () => {
		try {
			await fetch();
			isLoaded = true;
		} catch (error) {
			console.error('Error fetching GPU relations:', error);
		}
	});
</script>

<span class="flex items-center">
	{#if $nodes.length > 0 && isModelEnabled}
		<Drawer.Root
			bind:open
			onOpenChange={async (isOpen) => {
				if (isOpen) {
					if (!isLoaded) {
						await fetch();
						isLoaded = true;
					}
				} else {
					isLoaded = false;
				}
			}}
		>
			<Drawer.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon-sm' })}>
				<Icon icon="ph:graph" />
			</Drawer.Trigger>
			<Drawer.Content class="h-[77vh]">
				{#if isLoaded}
					<Drawer.Header>
						<Drawer.Title>{m.details()}</Drawer.Title>
						<Drawer.Description>
							<p>{m.gpu_relation_description()}</p>
							<p>{m.gpu_relation_guide_description()}</p>
						</Drawer.Description>
					</Drawer.Header>
					{#if open}
						<ComplexFlow.Flow initialNodes={$nodes} initialEdges={$edges} />
					{/if}
				{/if}
			</Drawer.Content>
		</Drawer.Root>
	{/if}
	{#if machineScope}
		<HoverCard.Root>
			<HoverCard.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon-sm' })}>
				<Icon icon="ph:graphics-card" />
			</HoverCard.Trigger>
			<HoverCard.Content class="m-0 h-fit w-fit rounded-lg p-0">
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
	<p class="p-2">{machine.gpuDevices.length}</p>
</span>
