<script lang="ts" module>
	import { type NodeProps } from '@xyflow/svelte';

	import Trigger from '../units/gpu-complex-trigger.svelte';
	import Details from '../units/gpu-details.svelte';

	import type { GPURelation_GPU, GPURelation_Pod_Device } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition,
		...restProps
	}: Omit<NodeProps, 'data'> & { data: { gpu: GPURelation_GPU; devices: GPURelation_Pod_Device[] } } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<Trigger data={data.gpu} {selected} {targetPosition} {sourcePosition} {...restProps} />
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<Details {data} />
	</HoverCard.Content>
</HoverCard.Root>
