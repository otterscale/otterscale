<script lang="ts" module>
	import { type NodeProps } from '@xyflow/svelte';

	import type { GPURelation_Pod } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';

	import Details from '../units/model-details.svelte';
	import Trigger from '../units/model-simple-trigger.svelte';
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
