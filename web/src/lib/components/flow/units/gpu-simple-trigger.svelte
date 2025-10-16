<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { GPURelation_GPU } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { dynamicPaths } from '$lib/path';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: GPURelation_GPU } =
		$props();

	const link = `${dynamicPaths.setupScopeKubernetes(page.params.scope).url}`;
</script>

<div
	class={cn(
		'bg-card relative rounded-full border p-4 shadow',
		selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
	)}
>
	<Icon icon="ph:graphics-card" class="size-10" />
	<p
		class="text-muted-foreground absolute bottom-0 left-1/2 max-w-[100px] -translate-x-1/2 truncate text-center text-xs whitespace-nowrap"
	>
		{data.type}
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
		'bg-card hover:bg-muted absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border p-2 shadow hover:cursor-default',
		selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
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
