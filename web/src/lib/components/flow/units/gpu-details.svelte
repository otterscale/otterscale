<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { GPURelation_GPU, GPURelation_Pod_Device } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { data }: { data: { gpu: GPURelation_GPU; devices: GPURelation_Pod_Device[] } } = $props();

	const usedCores = $derived(data.devices.reduce((a, device) => a + device.usedCores, 0));
	const vRAM = $derived(Number(data.gpu.memoryBytes));
	const usedvRAM = $derived(data.devices.reduce((a, device) => a + Number(device.usedMemoryBytes), 0));
	const usedCount = data.devices.length;

	const formatedVRAM = $derived(formatCapacity(vRAM));
	const formatedUsedVRAM = $derived(formatCapacity(usedvRAM));
</script>

<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
	<p class="text-lg font-bold">{data.gpu.type}</p>
	<div class="grid auto-rows-auto grid-cols-2 gap-1">
		<span class="col-span-2 flex items-center gap-2">
			<Icon icon="ph:identification-badge" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.id()}</p>
				{data.gpu.id}
			</div>
		</span>
		<span class="flex items-center gap-2">
			<Icon icon="ph:tag" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.index()}</p>
				{data.gpu.index}
			</div>
		</span>
		<span class="flex items-center gap-2">
			<Icon icon="ph:heartbeat" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.state()}</p>
				{#if data.gpu.health}
					<p class="text-sm text-green-500">{m.healthy()}</p>
				{:else}
					<p class="text-sm text-red-500">{m.unhealthy()}</p>
				{/if}
			</div>
		</span>
		<span class="flex items-center gap-2">
			<Icon icon="ph:cpu" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.cores()}</p>
				<div class="flex gap-1 font-mono">
					<p>{Number((usedCores * 100) / data.gpu.cores).toFixed(2)}%</p>
					<p class="text-muted-foreground">
						({usedCores}/{data.gpu.cores})
					</p>
				</div>
			</div>
		</span>
		<span class="flex items-center gap-2">
			<Icon icon="ph:memory" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.vram()}</p>
				<div class="flex gap-1 font-mono">
					<p>{Number((usedvRAM * 100) / vRAM).toFixed(2)}%</p>
					<p class="text-muted-foreground">
						({formatedUsedVRAM.value}{formatedUsedVRAM.unit}/{formatedVRAM.value}{formatedVRAM.unit})
					</p>
				</div>
			</div>
		</span>
		<span class="flex items-center gap-2">
			<Icon icon="ph:cube" class="text-muted-foreground size-6" />
			<div class="flex flex-col gap-0">
				<p class="text-muted-foreground text-sm">{m.gpu_count()}</p>
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
