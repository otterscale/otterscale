<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { GPURelation_Pod } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		data,
		selected,
		targetPosition,
		sourcePosition
	}: Omit<NodeProps, 'data'> & { data: { scope: string; pod: GPURelation_Pod } } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'relative rounded-full border bg-card p-4 shadow',
				selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
			)}
		>
			<Icon icon="ph:robot" class="size-10" />
			<p
				class="absolute bottom-0 left-1/2 max-w-[100px] -translate-x-1/2 truncate text-center text-xs whitespace-nowrap text-muted-foreground"
			>
				{data.pod.modelName}
			</p>
		</div>
		{#if targetPosition}
			<Handle type="target" position={targetPosition} class="invisible" />
		{/if}
		{#if sourcePosition}
			<Handle type="source" position={sourcePosition} class="invisible" />
		{/if}
		<div
			class={cn(
				'absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border bg-card p-2 shadow hover:cursor-default hover:bg-muted',
				selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
			)}
		>
			<Icon
				icon="ph:arrow-square-out"
				onclick={(e) => {
					e.stopPropagation();
					goto(
						resolve('/(auth)/scope/[scope]/applications/workloads', {
							scope: data.scope
						})
					);
				}}
			/>
		</div>
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
			<p class="text-lg font-bold">{data.pod.modelName}</p>
			<span class="flex items-center gap-2">
				<Icon icon="ph:cube" class="size-6 text-muted-foreground" />
				<div class="flex flex-col gap-0">
					<p class="text-sm text-muted-foreground">{m.namespace()}</p>
					{data.pod.namespace}
				</div>
			</span>
			<span class="flex items-center gap-2">
				<Icon icon="ph:cube" class="size-6 text-muted-foreground" />
				<div class="flex flex-col gap-0">
					<p class="text-sm text-muted-foreground">{m.application()}</p>
					<p>
						{data.pod.name}
					</p>
				</div>
			</span>
			<div class="flex flex-col gap-4">
				{#each data.pod.devices as device (device.gpuId)}
					{@const { value: usedMemoryValue, unit: usedMemoryUnit } = formatCapacity(
						Number(device.usedMemoryBytes)
					)}
					<div
						class="grid auto-rows-auto grid-cols-3 gap-1 rounded-lg bg-muted/50 p-4 hover:bg-muted"
					>
						<span class="col-span-3 flex items-center gap-2">
							<Icon icon="ph:identification-badge" class="size-6 text-muted-foreground" />
							<p class="flex flex-col gap-0">
								{device.gpuId}
							</p>
						</span>
						<span class="flex items-center gap-2">
							<Icon icon="ph:cpu" class="size-6 text-muted-foreground" />
							<div class="flex flex-col gap-0">
								<p class="text-sm text-muted-foreground">{m.cores()}</p>
								{device.usedCores}%
							</div>
						</span>
						<span class="flex items-center gap-2">
							<Icon icon="ph:memory" class="size-6 text-muted-foreground" />
							<div class="flex flex-col gap-0">
								<p class="text-sm text-muted-foreground">{m.vram()}</p>
								<p>
									{usedMemoryValue}
									{usedMemoryUnit}
								</p>
							</div>
						</span>
						<span class="flex items-center gap-2">
							<Icon icon="ph:cube" class="size-6 text-muted-foreground" />
							<div class="flex flex-col gap-0">
								<p class="text-sm text-muted-foreground">{m.gpu_count()}</p>
							</div>
						</span>
					</div>
				{/each}
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
