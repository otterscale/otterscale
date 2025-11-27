<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type {
		GPURelation_GPU,
		GPURelation_Machine
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
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
	}: Omit<NodeProps, 'data'> & { data: { machine: GPURelation_Machine; gpus: GPURelation_GPU[] } } =
		$props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'relative flex h-[150px] w-[300px] rounded-lg border bg-card p-2 hover:shadow',
				selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
			)}
		>
			<div
				class={cn(
					'absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border bg-card p-2 shadow hover:cursor-default hover:bg-muted',
					selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
				)}
			>
				<Icon
					icon="simple-icons:maas"
					class="size-5"
					onclick={(e) => {
						e.stopPropagation();
						goto(resolve('/(auth)/machines/metal/[id]', { id: data.machine.id }));
					}}
				/>
			</div>
			<div class="flex items-center justify-center p-4">
				<div class="flex items-center gap-2">
					<div class="size-fit rounded-full bg-muted-foreground/50 p-2">
						<Icon icon="ph:computer-tower" class="size-5" />
					</div>
					<div class="flex flex-col items-start justify-start">
						<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">
							{data.machine.hostname}
						</p>
						<p
							class="max-w-[200px] truncate text-xs font-light text-nowrap whitespace-nowrap text-muted-foreground"
						>
							{data.machine.id}
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
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
			<p class="text-lg font-bold">{data.machine.hostname}</p>
			<span class="flex items-center gap-2">
				<Icon icon="ph:identification-badge" class="size-6 text-muted-foreground" />
				<div class="flex flex-col gap-0">
					<p class="text-sm text-muted-foreground">{m.id()}</p>
					{data.machine.id}
				</div>
			</span>
			<div class="flex flex-col gap-4">
				{#each data.gpus as gpu (gpu.id)}
					{@const { value: usedMemoryValue, unit: usedMemoryUnit } = formatCapacity(
						Number(gpu.memoryBytes)
					)}
					<div
						class="grid auto-rows-auto grid-cols-3 gap-1 rounded-lg bg-muted/50 p-4 hover:bg-muted"
					>
						<div class="col-span-3">
							<p class="text-lg font-bold">{gpu.type}</p>
							<p class="text-muted-foreground">{gpu.id}</p>
						</div>
						<span class="flex items-center gap-2">
							<Icon icon="ph:cpu" class="size-6 text-muted-foreground" />
							<div class="flex flex-col gap-0">
								<p class="text-sm text-muted-foreground">{m.cores()}</p>
								{gpu.cores}%
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
								{gpu.count}
							</div>
						</span>
					</div>
				{/each}
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
