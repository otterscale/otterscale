<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import * as HoverCard from '$lib/components/ui/hover-card';
	import { cn } from '$lib/utils';

	export type DataType = {
		name: string;
		model: string;
		icon: string;
	};
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: DataType } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'bg-card rounded-full border p-4 shadow',
				selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
			)}
		>
			<Icon icon="ph:graphics-card-duotone" class="size-10" />
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
			<div class="bg-muted-foreground/50 size-fit rounded-full">
				<Icon icon={data.icon} class="size-5" />
			</div>
			<div>
				<p class="text-base text-nowrap whitespace-nowrap">{data.name}</p>
				<p class="text-muted-foreground text-xs font-light text-nowrap whitespace-nowrap">
					{data.model}
				</p>
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
