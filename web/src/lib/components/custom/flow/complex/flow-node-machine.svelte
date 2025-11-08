<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { cn } from '$lib/utils';

	export type DataType = {
		name: string;
		ip: string;
		icon: string;
	};
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition
	}: Omit<NodeProps, 'data'> & { data: DataType } = $props();
</script>

<div
	class={cn(
		'relative flex h-[150px] w-[300px] rounded-lg border bg-card p-2 hover:shadow',
		selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
	)}
>
	<div
		class={cn(
			'absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border bg-card p-2 shadow',
			selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
		)}
	>
		<Icon icon="ph:computer-tower-duotone" class="size-5" />
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="flex gap-2">
			<div class="size-fit rounded-full bg-muted-foreground/50 p-1">
				<Icon icon={data.icon} class="size-5" />
			</div>
			<div>
				<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">{data.name}</p>
				<p
					class="max-w-[200px] truncate text-xs font-light text-nowrap whitespace-nowrap text-muted-foreground"
				>
					{data.ip}
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
