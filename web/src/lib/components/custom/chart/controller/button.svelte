<script lang="ts">
	let {
		activeChart = $bindable(),
		keys,
		data,
	}: {
		activeChart: string;
		keys: string[];
		data: Record<string, number>;
	} = $props();
</script>

<div class="flex">
	{#each keys as key (key)}
		<button
			data-active={activeChart === key}
			class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-6 sm:py-4"
			onclick={() => (activeChart = key)}
		>
			<span class="text-muted-foreground text-xs">
				{key}
			</span>
			<span class="text-lg leading-none font-bold sm:text-3xl">
				{data[key].toLocaleString?.() ?? NaN}
			</span>
		</button>
	{/each}
</div>
