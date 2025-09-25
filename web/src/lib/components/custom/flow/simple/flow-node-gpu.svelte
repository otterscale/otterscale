<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import type { GpuInfo } from '$lib/api/essential/v1/essential_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { formatCapacity } from '$lib/formatter';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: GpuInfo } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'bg-card rounded-full border p-4 shadow',
				selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
			)}
		>
			<Icon icon="ph:graphics-card" class="size-10" />
		</div>
		{#if targetPosition}
			<Handle type="target" position={targetPosition} class="invisible" />
		{/if}
		{#if sourcePosition}
			<Handle type="source" position={sourcePosition} class="invisible" />
		{/if}
	</HoverCard.Trigger>
	<HoverCard.Content class="w-fit">
		<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
			<p class="text-lg font-bold">{data.physicalGpuUuid}</p>
			<div class="grid auto-rows-auto grid-cols-2">
				<span class="flex items-center gap-2">
					<Icon icon="ph:cpu" class="text-muted-foreground size-6" />
					<div class="flex flex-col gap-0">
						<p class="text-muted-foreground text-sm">cores</p>
						{data.vcoresPercent}
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:memory" class="text-muted-foreground size-6" />
					<div class="flex flex-col gap-0">
						<p class="text-muted-foreground text-sm">vRAM</p>
						<p>
							{formatCapacity(Number(data.vramMib) * 1024 * 1024).value}
							{formatCapacity(Number(data.vramMib) * 1024 * 1024).unit}
						</p>
					</div>
				</span>
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
