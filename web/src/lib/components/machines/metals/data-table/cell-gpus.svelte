<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Simple } from '$lib/components/custom/flow/index';
	import { traverse } from '$lib/components/custom/flow/utils.svelte';
	import * as Table from '$lib/components/custom/table';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };

	const nodes: Node[] = [
		{
			id: 'machine',
			type: 'machine',
			data: {
				name: 'otterscale-vm141.maas',
				ip: '10.102.197.141',
				icon: 'simple-icons:maas',
			},
			position,
		},
		{
			id: 'gpu1',
			type: 'gpu',
			data: {
				name: 'NVIDIA RTX 4090',
				model: 'NVIDIA',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'gpu2',
			type: 'gpu',
			data: {
				name: 'AMD Radeon RX 7900 XTX',
				model: 'AMD',
				icon: 'simple-icons:amd',
			},
			position,
		},
		{
			id: 'gpu3',
			type: 'gpu',
			data: {
				name: 'NVIDIA A100',
				model: 'NVIDIA',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'gpu4',
			type: 'gpu',
			data: {
				name: 'NVIDIA RTX 3080',
				model: 'NVIDIA',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'gpu5',
			type: 'gpu',
			data: {
				name: 'AMD Radeon RX 6800 XT',
				model: 'AMD',
				icon: 'simple-icons:amd',
			},
			position,
		},
		{
			id: 'gpu6',
			type: 'gpu',
			data: {
				name: 'NVIDIA Tesla V100',
				model: 'NVIDIA',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'gpu7',
			type: 'gpu',
			data: {
				name: 'Intel Arc A770',
				model: 'Intel',
				icon: 'simple-icons:intel',
			},
			position,
		},
		{
			id: 'gpu8',
			type: 'gpu',
			data: {
				name: 'NVIDIA Quadro RTX 6000',
				model: 'NVIDIA',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'vgpu1',
			type: 'vgpu',
			data: {
				name: 'NVIDIA0',
				physical: 'NVIDIA Quadro RTX 6000',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'vgpu2',
			type: 'vgpu',
			data: {
				name: 'NVIDIA1',
				physical: 'NVIDIA Quadro RTX 6000',
				icon: 'simple-icons:nvidia',
			},
			position,
		},
		{
			id: 'model1',
			type: 'model',
			data: {
				name: 'gpt-oss9',
				framework: 'vllm',
				icon: 'simple-icons:openai',
			},
			position,
		},
		{
			id: 'model2',
			type: 'model',
			data: {
				name: 'llama-3-70b',
				framework: 'transformers',
				icon: 'simple-icons:meta',
			},
			position,
		},
		{
			id: 'model3',
			type: 'model',
			data: {
				name: 'mixtral-8x7b',
				framework: 'vllm',
				icon: 'simple-icons:mistral',
			},
			position,
		},
	];

	const edges: Edge[] = [
		{ id: '1', type: 'edge', source: 'model1', target: 'gpu1', animated: true, selectable: false },
		{ id: '2', type: 'edge', source: 'model1', target: 'gpu2', animated: true, selectable: false },
		{ id: '3', type: 'edge', source: 'model2', target: 'gpu3', animated: true, selectable: false },
		{ id: '4', type: 'edge', source: 'model2', target: 'gpu4', animated: true, selectable: false },
		{ id: '5', type: 'edge', source: 'model2', target: 'gpu5', animated: true, selectable: false },
		{ id: '6', type: 'edge', source: 'model3', target: 'gpu6', animated: true, selectable: false },
		{ id: '7', type: 'edge', source: 'model3', target: 'gpu7', animated: true, selectable: false },
		{ id: '8', type: 'edge', source: 'model3', target: 'vgpu1', animated: true, selectable: false },
		{ id: '19', type: 'edge', source: 'model3', target: 'vgpu2', animated: true, selectable: false },
		{ id: '17', type: 'edge', source: 'vgpu1', target: 'gpu8', animated: true, selectable: false },
		{ id: '18', type: 'edge', source: 'vgpu2', target: 'gpu8', animated: true, selectable: false },
		{ id: '9', type: 'edge', source: 'gpu1', target: 'machine', animated: true, selectable: false },
		{ id: '10', type: 'edge', source: 'gpu2', target: 'machine', animated: true, selectable: false },
		{ id: '11', type: 'edge', source: 'gpu3', target: 'machine', animated: true, selectable: false },
		{ id: '12', type: 'edge', source: 'gpu4', target: 'machine', animated: true, selectable: false },
		{ id: '13', type: 'edge', source: 'gpu5', target: 'machine', animated: true, selectable: false },
		{ id: '14', type: 'edge', source: 'gpu6', target: 'machine', animated: true, selectable: false },
		{ id: '15', type: 'edge', source: 'gpu7', target: 'machine', animated: true, selectable: false },
		{ id: '16', type: 'edge', source: 'gpu8', target: 'machine', animated: true, selectable: false },
	];

	let initialNodes: Node[] = $state([]);
	let initialEdges: Edge[] = $state([]);

	const selectedNode = nodes.find((node) => node.selected);
	if (selectedNode) {
		const { nodestoFocus, edgesToFocus } = traverse(edges, selectedNode.id);

		initialNodes = nodes.map((node) => ({
			...node,
			selected: nodestoFocus.has(node.id),
		}));

		initialEdges = edges.map((edge) => ({
			...edge,
			selected: edgesToFocus.has(edge.id),
		}));
	} else {
		initialNodes = nodes;
		initialEdges = edges;
	}

	let open = $state(false);

	let {
		machine,
	}: {
		machine: Machine;
	} = $props();
</script>

<Sheet.Root bind:open>
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
			<Simple.Flow {initialNodes} {initialEdges} />
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
								<Table.Cell
									>{gpuDevice.productName !== ''
										? gpuDevice.productName
										: gpuDevice.productId}</Table.Cell
								>
								<Table.Cell
									>{gpuDevice.vendorName !== ''
										? gpuDevice.vendorName
										: gpuDevice.vendorId}</Table.Cell
								>
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
