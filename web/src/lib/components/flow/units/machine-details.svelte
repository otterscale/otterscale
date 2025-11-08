<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type {
		GPURelation_GPU,
		GPURelation_Machine
	} from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { data }: { data: { machine: GPURelation_Machine; gpus: GPURelation_GPU[] } } = $props();
</script>

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
		{#each data.gpus as gpu}
			{@const { value: usedMemoryValue, unit: usedMemoryUnit } = formatCapacity(
				Number(gpu.memoryBytes)
			)}
			<div class="grid auto-rows-auto grid-cols-3 gap-1 rounded-lg bg-muted/50 p-4 hover:bg-muted">
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
