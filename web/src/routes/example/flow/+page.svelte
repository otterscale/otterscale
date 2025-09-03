<script lang="ts">
	import { Root } from '$lib/components/custom/flow/index';
	import type { GPU, Machine, Model } from '$lib/components/custom/flow/types';
	import { type Edge, type Node } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };

	export const initialNodes: Node[] = [
		{
			id: 'machine1',
			type: 'machine',
			data: {
				typeIcon: 'ph:computer-tower',
				name: 'juju-vm.maas',
				ip: '10.102.197.19',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'machine2',
			type: 'machine',
			data: {
				typeIcon: 'ph:computer-tower',
				name: 'otterscale-vm141.maas',
				ip: '10.102.197.141',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'machine3',
			type: 'machine',
			data: {
				typeIcon: 'ph:computer-tower',
				name: 'development-vm183.maas',
				ip: '10.102.197.164',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'machine4',
			type: 'machine',
			data: {
				typeIcon: 'ph:computer-tower',
				name: 'development-vm185.maas',
				icon: 'simple-icons:maas',
			} as Machine,
			position,
		},
		{
			id: 'gpu1',
			type: 'gpu',
			data: {
				typeIcon: 'ph:graphics-card',
				name: 'g4201-6000adax8-197-110',
				manufacturer: 'nvidia0',
				icon: 'simple-icons:nvidia',
			} as GPU,
			position,
		},
		{
			id: 'gpu2',
			type: 'gpu',
			data: {
				typeIcon: 'ph:graphics-card',
				name: 'g4201-a6000x8-197-109',
				manufacturer: 'nvidia1',
				icon: 'simple-icons:nvidia',
			} as GPU,
			position,
		},
		{
			id: 'gpu3',
			type: 'gpu',
			data: {
				typeIcon: 'ph:graphics-card',
				name: 'proxmox-4090x2-197-114',
				manufacturer: 'nvidia0',
				icon: 'simple-icons:nvidia',
			} as GPU,
			position,
		},
		{
			id: 'model1',
			type: 'model',
			data: {
				typeIcon: 'ph:robot',
				name: 'gpt-oss9',
				framework: 'vllm',
				icon: 'simple-icons:openai',
			} as Model,
			position,
		},
		{
			id: 'model2',
			type: 'model',
			data: {
				typeIcon: 'ph:robot',
				name: 'qwen3',
				framework: 'vllm',
				icon: 'simple-icons:alibabacloud',
			} as Model,
			position,
		},
	];

	export const initialEdges: Edge[] = [
		{ id: 'edge-machine2-gpu1', type: 'edge', source: 'machine2', target: 'gpu1', animated: true },
		{ id: 'edge-machine3-gpu2', type: 'edge', source: 'machine3', target: 'gpu2', animated: true },
		{ id: 'edge-machine3-gpu3', type: 'edge', source: 'machine3', target: 'gpu3', animated: true },
		{ id: 'edge-machine4-gpu3', type: 'edge', source: 'machine4', target: 'gpu3', animated: true },
		{ id: 'edge-gpu1-model1', type: 'edge', source: 'gpu1', target: 'model1', animated: true },
		{ id: 'edge-gpu2-model1', type: 'edge', source: 'gpu1', target: 'model2', animated: true },
		{ id: 'edge-gpu2-model2', type: 'edge', source: 'gpu2', target: 'model2', animated: true },
	];
</script>

<Root {initialNodes} {initialEdges} />
