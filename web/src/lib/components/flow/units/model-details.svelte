<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { GPURelation_Pod } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { data }: { data: GPURelation_Pod } = $props();
</script>

<div class="flex flex-col gap-2 p-2 text-base text-nowrap whitespace-nowrap">
	<p class="text-lg font-bold">{data.modelName}</p>
	<span class="flex items-center gap-2">
		<Icon icon="ph:cube" class="text-muted-foreground size-6" />
		<div class="flex flex-col gap-0">
			<p class="text-muted-foreground text-sm">{m.namespace()}</p>
			{data.namespace}
		</div>
	</span>
	<span class="flex items-center gap-2">
		<Icon icon="ph:cube" class="text-muted-foreground size-6" />
		<div class="flex flex-col gap-0">
			<p class="text-muted-foreground text-sm">{m.application()}</p>
			<p>
				{data.name}
			</p>
		</div>
	</span>
	<div class="flex flex-col gap-4">
		{#each data.devices as device}
			{@const { value: usedMemoryValue, unit: usedMemoryUnit } = formatCapacity(Number(device.usedMemoryBytes))}
			<div class="bg-muted/50 hover:bg-muted grid auto-rows-auto grid-cols-3 gap-1 rounded-lg p-4">
				<span class="col-span-3 flex items-center gap-2">
					<Icon icon="ph:identification-badge" class="text-muted-foreground size-6" />
					<p class="flex flex-col gap-0">
						{device.gpuId}
					</p>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:cpu" class="text-muted-foreground size-6" />
					<div class="flex flex-col gap-0">
						<p class="text-muted-foreground text-sm">{m.cores()}</p>
						{device.usedCores}%
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:memory" class="text-muted-foreground size-6" />
					<div class="flex flex-col gap-0">
						<p class="text-muted-foreground text-sm">{m.vram()}</p>
						<p>
							{usedMemoryValue}
							{usedMemoryUnit}
						</p>
					</div>
				</span>
				<span class="flex items-center gap-2">
					<Icon icon="ph:cube" class="text-muted-foreground size-6" />
					<div class="flex flex-col gap-0">
						<p class="text-muted-foreground text-sm">{m.gpu_count()}</p>
						1
					</div>
				</span>
			</div>
		{/each}
	</div>
</div>
