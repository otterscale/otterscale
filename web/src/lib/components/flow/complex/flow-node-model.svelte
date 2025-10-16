<script lang="ts" module>
	import { type NodeProps } from '@xyflow/svelte';

	import Trigger from '../units/model-complex-trigger.svelte';
	import Details from '../units/model-details.svelte';

	import type { GPURelation_Pod } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition,
		...restProps
	}: Omit<NodeProps, 'data'> & { data: GPURelation_Pod } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<Trigger {data} {selected} {targetPosition} {sourcePosition} {...restProps} />
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<Details {data} />
	</HoverCard.Content>
</HoverCard.Root>
