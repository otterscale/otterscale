<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { GPURelation_Machine } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition
	}: Omit<NodeProps, 'data'> & { data: GPURelation_Machine } = $props();

	const link = resolve('/(auth)/machines/metal/[id]', { id: data.id });
</script>

<div
	class={cn(
		'relative flex h-[150px] w-[300px] rounded-lg border bg-card p-2 hover:shadow',
		selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
	)}
>
	<div
		class={cn(
			'absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border bg-card p-2 shadow hover:cursor-default hover:bg-muted',
			selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
		)}
	>
		<Icon
			icon="simple-icons:maas"
			class="size-5"
			onclick={(e) => {
				e.stopPropagation();
				goto(link);
			}}
		/>
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="flex items-center gap-2">
			<div class="size-fit rounded-full bg-muted-foreground/50 p-2">
				<Icon icon="ph:computer-tower" class="size-5" />
			</div>
			<div class="justufy-start flex flex-col items-start">
				<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">
					{data.hostname}
				</p>
				<p
					class="max-w-[200px] truncate text-xs font-light text-nowrap whitespace-nowrap text-muted-foreground"
				>
					{data.id}
				</p>
			</div>
		</div>
	</div>
	{#if targetPosition}
		<Handle type="target" position={targetPosition} class="invisible" />
	{/if}
	{#if sourcePosition}
		<Handle type="source" position={sourcePosition} class="invisible" />
	{/if}
</div>
