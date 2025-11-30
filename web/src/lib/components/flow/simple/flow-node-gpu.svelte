<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type {
		GPURelation_GPU,
		GPURelation_Pod_Device
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
	}: Omit<NodeProps, 'data'> & {
		data: { scope: string; gpu: GPURelation_GPU; devices: GPURelation_Pod_Device[] };
	} = $props();

	const usedCores = $derived(data.devices.reduce((a, device) => a + device.usedCores, 0));
	const vRAM = $derived(Number(data.gpu.memoryBytes));
	const usedvRAM = $derived(
		data.devices.reduce((a, device) => a + Number(device.usedMemoryBytes), 0)
	);
	const usedCount = data.devices.length;

	const formatedVRAM = $derived(formatCapacity(vRAM));
	const formatedUsedVRAM = $derived(formatCapacity(usedvRAM));
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<div
			class={cn(
				'relative rounded-full border bg-card p-4 shadow',
				selected ? 'bg-primary-foreground ring-1 ring-primary' : 'bg-card ring-0'
			)}
		>
			<Icon icon="ph:graphics-card" class="size-10" />
			<p
				class="absolute bottom-0 left-1/2 max-w-[100px] -translate-x-1/2 truncate text-center text-xs whitespace-nowrap text-muted-foreground"
			>
				{data.gpu.type}
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
					goto(resolve('/(auth)/scope/[scope]/setup/kubernetes', { scope: data.scope }));
				}}
			/>
		</div>
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
			<p class="text-lg font-bold">{data.gpu.type}</p>
			<div class="grid auto-rows-auto grid-cols-2 gap-1">
				<span class="col-span-2 flex items-center gap-2">
					<Icon icon="ph:identification-badge" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.id()}</p>
						{data.gpu.id}
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:tag" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.index()}</p>
						{data.gpu.index}
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:heartbeat" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.state()}</p>
						{#if data.gpu.health}
							<p class="text-sm text-green-500">{m.healthy()}</p>
						{:else}
							<p class="text-sm text-red-500">{m.unhealthy()}</p>
						{/if}
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:cpu" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.cores()}</p>
						<div class="flex gap-1 font-mono">
							<p>{Number((usedCores * 100) / data.gpu.cores).toFixed(2)}%</p>
							<p class="text-muted-foreground">
								({usedCores}/{data.gpu.cores})
							</p>
						</div>
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:memory" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.vram()}</p>
						<div class="flex gap-1 font-mono">
							<p>{Number((usedvRAM * 100) / vRAM).toFixed(2)}%</p>
							<p class="text-muted-foreground">
								({formatedUsedVRAM.value}{formatedUsedVRAM.unit}/{formatedVRAM.value}{formatedVRAM.unit})
							</p>
						</div>
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:cube" class="size-6 text-muted-foreground" />
					<div class="flex flex-col gap-0">
						<p class="text-sm text-muted-foreground">{m.gpu_count()}</p>
						<div class="flex gap-1 font-mono">
							<p>{Number((usedCount * 100) / data.gpu.count).toFixed(2)}%</p>
							<p class="text-muted-foreground">
								({usedCount}/{data.gpu.count})
							</p>
						</div>
					</div>
				</span>
			</div>
		</div>
	</HoverCard.Content>
</HoverCard.Root>
