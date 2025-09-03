<script lang="ts">
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	export type DataType = {
		typeIcon: string;
		name: string;
		manufacturer: string;
		icon: string;
	};

	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: DataType } = $props();
</script>

<div
	class={cn(
		'relative flex min-h-[100px] min-w-[200px] rounded-lg border p-2 hover:shadow',
		selected ? 'bg-primary-foreground ring-primary animate-pulse ring-[0.5px]' : 'bg-none ring-0',
	)}
>
	<div
		class={cn(
			'bg-card absolute right-1 top-1 -translate-y-1/2 translate-x-1/2 rounded-full border p-2 shadow',
			selected ? 'bg-primary-foreground ring-primary ring-[0.5px]' : 'bg-none ring-0',
		)}
	>
		<Icon icon={data.typeIcon} class="size-5" />
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="grid grid-cols-[auto_1fr] items-start">
			<div class="bg-muted-foreground/50 rounded-full p-2">
				<Icon icon={data.icon} class="size-5 self-start" />
			</div>
			<div class="p-1.5">
				<p class="text-base">{data.name}</p>
				<p class="text-muted-foreground text-xs font-light">{data.manufacturer}</p>
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
