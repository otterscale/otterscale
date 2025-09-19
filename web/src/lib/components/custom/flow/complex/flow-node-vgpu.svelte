<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { cn } from '$lib/utils';

	export type DataType = {
		name: string;
		physical: string;
		icon: string;
	};
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: DataType } = $props();
</script>

<div
	class={cn(
		'bg-card relative flex h-[150px] w-[300px] rounded-lg border p-2 hover:shadow',
		selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
	)}
>
	<div
		class={cn(
			'bg-card absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border p-2 shadow',
			selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
		)}
	>
		<Icon icon="ph:graphics-card-thin" class="size-5" />
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="flex gap-2">
			<div class="bg-muted-foreground/50 size-fit rounded-full p-1">
				<Icon icon={data.icon} class="size-5" />
			</div>
			<div>
				<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">{data.name}</p>
				<p
					class="text-muted-foreground max-w-[200px] truncate text-xs font-light text-nowrap whitespace-nowrap"
				>
					{data.physical}
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
