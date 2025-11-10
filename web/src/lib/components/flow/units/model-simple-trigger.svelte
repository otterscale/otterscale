<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { GPURelation_Pod } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition
	}: Omit<NodeProps, 'data'> & { data: GPURelation_Pod & { scope: string } } = $props();

	const link = resolve('/(auth)/scope/[scope]/applications/workloads', {
		scope: data.scope
	});
</script>

<div
	class={cn(
		'relative rounded-full border bg-card p-4 shadow',
		selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
	)}
>
	<Icon icon="ph:robot" class="size-10" />
	<p
		class="absolute bottom-0 left-1/2 max-w-[100px] -translate-x-1/2 truncate text-center text-xs whitespace-nowrap text-muted-foreground"
	>
		{data.modelName}
	</p>
</div>
{#if targetPosition}
	<Handle type="target" position={targetPosition} class="invisible" />
{/if}
{#if sourcePosition}
	<Handle type="source" position={sourcePosition} class="invisible" />
{/if}
<div
	class={cn(
		'absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border bg-card p-2 shadow hover:cursor-default hover:bg-muted',
		selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
	)}
>
	<Icon
		icon="ph:arrow-square-out"
		onclick={(e) => {
			e.stopPropagation();
			goto(link);
		}}
	/>
</div>
