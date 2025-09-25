<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import type { PodInfo } from '$lib/api/essential/v1/essential_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { cn } from '$lib/utils';

	export type DataType = {
		name: string;
		framework: string;
		icon: string;
	};
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: PodInfo } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'bg-card rounded-full border p-4 shadow',
				selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
			)}
		>
			<Icon icon="ph:robot" class="size-10" />
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
			<p class="text-lg font-bold">{data.modelName}</p>
			<span class="flex items-center gap-2">
				<Icon icon="ph:cube" class="text-muted-foreground size-6" />
				<div class="flex flex-col gap-0">
					<p class="text-muted-foreground text-sm">namespace</p>
					{data.namespace}
				</div>
			</span>
			<span class="flex items-center gap-2">
				<Icon icon="ph:cube" class="text-muted-foreground size-6" />
				<div class="flex flex-col gap-0">
					<p class="text-muted-foreground text-sm">application</p>
					<p>
						{data.name}
					</p>
				</div>
			</span>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
