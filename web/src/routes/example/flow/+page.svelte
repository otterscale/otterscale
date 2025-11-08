<script lang="ts">
	import { type Edge, type Node } from '@xyflow/svelte';

	import {
		Complex,
		type GPU,
		type Machine,
		type Model,
		type VGPU
	} from '$lib/components/custom/flow/index';
	import '@xyflow/svelte/dist/style.css';

	const position = { x: 0, y: 0 };

	const initialNodes: Node[] = [
		{
			id: 'machine1',
			type: 'machine',
			data: {
				name: 'otterscale-vm141.maas',
				ip: '10.102.197.141',
				icon: 'simple-icons:maas'
			} as Machine,
			position
		},
		{
			id: 'machine2',
			type: 'machine',
			data: {
				name: 'development-vm183.maas',
				ip: '10.102.197.164',
				icon: 'simple-icons:maas'
			} as Machine,
			position
		},
		{
			id: 'gpu2-vgpu1',
			type: 'vgpu',
			data: {
				name: 'nvidia0',
				physical: 'proxmox-4090x2-197-113',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},
		{
			id: 'gpu2-vgpu2',
			type: 'vgpu',
			data: {
				name: 'nvidia1',
				physical: 'proxmox-4090x2-197-113',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},
		{
			id: 'gpu1-vgpu1',
			type: 'vgpu',
			data: {
				name: 'nvidia0',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},
		{
			id: 'gpu1-vgpu2',
			type: 'vgpu',
			data: {
				name: 'nvidia1',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},
		{
			id: 'gpu1-vgpu3',
			type: 'vgpu',
			data: {
				name: 'nvidia2',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},
		{
			id: 'gpu1-vgpu4',
			type: 'vgpu',
			data: {
				name: 'nvidia3',
				physical: 'g4201-6000adax8-197-110',
				icon: 'simple-icons:nvidia'
			} as VGPU,
			position
		},

		{
			id: 'gpu1',
			type: 'gpu',
			data: {
				name: 'g4201-6000adax8-197-110',
				model: 'NVIDIA RTX 6000 Ada Generation',
				icon: 'simple-icons:nvidia'
			} as GPU,
			position
		},
		{
			id: 'gpu2',
			type: 'gpu',
			data: {
				name: 'proxmox-4090x2-197-113',
				model: 'NVIDIA GeForce RTX 4090',
				icon: 'simple-icons:nvidia'
			} as GPU,
			position
		},
		{
			id: 'model',
			type: 'model',
			data: {
				name: 'gpt-oss9',
				framework: 'vllm',
				icon: 'simple-icons:openai'
			} as Model,
			position
		}
	];

	const initialEdges: Edge[] = [
		{ id: '1', type: 'edge', source: 'gpu1-vgpu1', target: 'model', animated: true },
		{ id: '2', type: 'edge', source: 'gpu1-vgpu2', target: 'model', animated: true },
		{ id: '3', type: 'edge', source: 'gpu1-vgpu3', target: 'model', animated: true },
		{ id: '4', type: 'edge', source: 'gpu1-vgpu4', target: 'model', animated: true },
		{ id: '18', type: 'edge', source: 'gpu2-vgpu1', target: 'model', animated: true },
		{ id: '19', type: 'edge', source: 'gpu2-vgpu2', target: 'model', animated: true },
		{ id: '9', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu1', animated: true },
		{ id: '10', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu2', animated: true },
		{ id: '11', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu3', animated: true },
		{ id: '12', type: 'edge', source: 'gpu1', target: 'gpu1-vgpu4', animated: true },
		{ id: '20', type: 'edge', source: 'gpu2', target: 'gpu2-vgpu1', animated: true },
		{ id: '21', type: 'edge', source: 'gpu2', target: 'gpu2-vgpu2', animated: true },
		{ id: '17', type: 'edge', source: 'machine1', target: 'gpu1', animated: true },
		{ id: '22', type: 'edge', source: 'machine2', target: 'gpu2', animated: true }
	];
</script>

<div class="h-[1000px]">
	<Complex.Flow {initialNodes} {initialEdges} />
</div>
