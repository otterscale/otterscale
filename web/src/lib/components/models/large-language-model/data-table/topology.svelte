<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Edge, type Node } from '@xyflow/svelte';

	import { Root, type GPU, type Machine, type Model, type VGPU } from '$lib/components/custom/flow/index';
	import { traverse } from '$lib/components/custom/flow/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog/index.js';

	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };

	const nodes: Node[] = [
		{
			id: 'machine1',
			type: 'machine',
			data: {
				name: 'otterscale-vm141.maas',
				ip: '10.102.197.141',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'machine2',
			type: 'machine',
			data: {
				name: 'development-vm183.maas',
				ip: '10.102.197.164',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'gpu2-vgpu1',
			type: 'vgpu',
			data: {
				name: 'nvidia0',
				physical: 'proxmox-4090x2-197-113',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},
		{
			id: 'gpu2-vgpu2',
			type: 'vgpu',
			data: {
				name: 'nvidia1',
				physical: 'proxmox-4090x2-197-113',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},
		{
			id: 'gpu1-vgpu1',
			type: 'vgpu',
			data: {
				name: 'nvidia0',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},
		{
			id: 'gpu1-vgpu2',
			type: 'vgpu',
			data: {
				name: 'nvidia1',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},
		{
			id: 'gpu1-vgpu3',
			type: 'vgpu',
			data: {
				name: 'nvidia2',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},
		{
			id: 'gpu1-vgpu4',
			type: 'vgpu',
			data: {
				name: 'nvidia3',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia',
			} as VGPU,
			position,
		},

		{
			id: 'gpu1',
			type: 'gpu',
			data: {
				name: 'g4201-6000adax8-197-110',
				model: 'NVIDIA RTX 6000 Ada Generation',
				icon: 'simple-icons:nvidia',
			} as GPU,
			position,
		},
		{
			id: 'gpu2',
			type: 'gpu',
			data: {
				name: 'proxmox-4090x2-197-113',
				model: 'NVIDIA GeForce RTX 4090',
				icon: 'simple-icons:nvidia',
			} as GPU,
			position,
		},
		{
			id: 'model',
			type: 'model',
			data: {
				name: 'gpt-oss9',
				framework: 'vllm',
				icon: 'simple-icons:openai',
			} as Model,
			position,
			selected: true,
		},
	];

	const edges: Edge[] = [
		{ id: '1', type: 'edge', source: 'gpu1-vgpu1', target: 'model', animated: true, selectable: false },
		{ id: '2', type: 'edge', source: 'gpu1-vgpu2', target: 'model', animated: true, selectable: false },
		{ id: '3', type: 'edge', source: 'gpu1-vgpu3', target: 'model', animated: true, selectable: false },
		{ id: '4', type: 'edge', source: 'gpu1-vgpu4', target: 'model', animated: true, selectable: false },
		{ id: '18', type: 'edge', source: 'gpu2-vgpu1', target: 'model', animated: true, selectable: false },
		{ id: '19', type: 'edge', source: 'gpu2-vgpu2', target: 'model', animated: true, selectable: false },
		{ id: '9', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu1', animated: true, selectable: false },
		{ id: '10', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu2', animated: true, selectable: false },
		{ id: '11', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu3', animated: true, selectable: false },
		{ id: '12', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu4', animated: true, selectable: false },
		{ id: '20', type: 'edge', source: 'gpu2', target: 'gpu2-vgpu1', animated: true, selectable: false },
		{ id: '21', type: 'edge', source: 'gpu2', target: 'gpu2-vgpu2', animated: true, selectable: false },
		{ id: '17', type: 'edge', source: 'machine1', target: 'gpu1', animated: true, selectable: false },
		{ id: '22', type: 'edge', source: 'machine2', target: 'gpu2', animated: true, selectable: false },
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
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
		<Icon icon="ph:graph" />
	</Dialog.Trigger>
	{#if open}
		<Dialog.Content class="min-h-[77vh] min-w-[77vw]">
			<Root {initialNodes} {initialEdges} />
		</Dialog.Content>
	{/if}
</Dialog.Root>
