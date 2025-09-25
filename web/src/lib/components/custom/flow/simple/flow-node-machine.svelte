<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import type { PodInfo } from '$lib/api/essential/v1/essential_pb';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: PodInfo } = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger>
			<div
				class={cn(
					'bg-card rounded-full border p-4 shadow',
					selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
				)}
			>
				<Icon icon="ph:computer-tower" class="size-10" />
			</div>
			{#if targetPosition}
				<Handle type="target" position={targetPosition} class="invisible" />
			{/if}
			{#if sourcePosition}
				<Handle type="source" position={sourcePosition} class="invisible" />
			{/if}
		</Tooltip.Trigger>
		<Tooltip.Content>
			{data.machineName}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
