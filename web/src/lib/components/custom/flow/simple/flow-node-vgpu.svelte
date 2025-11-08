<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import * as HoverCard from '$lib/components/ui/hover-card';
	import { cn } from '$lib/utils';

	export type DataType = {
		name: string;
		physical: string;
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

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'rounded-full border bg-card p-4 shadow',
				selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
			)}
		>
			<Icon icon="ph:graphics-card-thin" class="size-10" />
		</div>
		{#if targetPosition}
			<Handle type="target" position={targetPosition} class="invisible" />
		{/if}
		{#if sourcePosition}
			<Handle type="source" position={sourcePosition} class="invisible" />
		{/if}
	</HoverCard.Trigger>
	<HoverCard.Content class="w-fit">
		<div class="flex gap-2 p-2">
			<div class="size-fit rounded-full bg-muted-foreground/50">
				<Icon icon={data.icon} class="size-5" />
			</div>
			<div>
				<p class="text-base text-nowrap whitespace-nowrap">{data.name}</p>
				<p class="text-xs font-light text-nowrap whitespace-nowrap text-muted-foreground">
					{data.physical}
				</p>
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
