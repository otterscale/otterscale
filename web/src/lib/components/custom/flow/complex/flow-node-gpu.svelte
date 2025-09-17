<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { type GpuInfo } from '$lib/api/essential/v1/essential_pb';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: GpuInfo } = $props();
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
		<Icon icon="simple-icons:kubernetes" class="size-5" />
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="flex gap-2">
			<div class="bg-muted-foreground/50 size-fit rounded-full p-2">
				<Icon icon="ph:graphics-card" class="size-5" />
			</div>
			<div class="justufy-start flex flex-col items-start">
				<Tooltip.Root>
					<Tooltip.Trigger>
						<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">
							{data.physicalGpuUuid}
						</p>
					</Tooltip.Trigger>
					<Tooltip.Content>
						{data.physicalGpuUuid}
					</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						<p
							class="text-muted-foreground max-w-[200px] truncate text-xs font-light text-nowrap whitespace-nowrap"
						>
							{data.vgpuCores} cores · {data.vgpuVram} ram
						</p>
					</Tooltip.Trigger>
					<Tooltip.Content>
						{data.vgpuCores} cores · {data.vgpuVram} ram
					</Tooltip.Content>
				</Tooltip.Root>
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
