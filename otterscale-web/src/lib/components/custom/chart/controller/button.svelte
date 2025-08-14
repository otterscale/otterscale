<script lang="ts">
    let {
        activeChart = $bindable(),
        chartConfig,
        total,
        options = ["desktop", "mobile", "notebook"]
    }: {
        activeChart: string;
        chartConfig: Record<string, { label: string; color?: string }>;
        total: Record<string, number>;
        options?: string[];
    } = $props();
</script>

<div class="flex">
    {#each options as key (key)}
        {@const chart = key as keyof typeof chartConfig}
        <button
            data-active={activeChart === chart}
            class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-l sm:border-t-0 sm:px-6 sm:py-4"
            onclick={() => (activeChart = chart)}
        >
        <span class="text-muted-foreground text-xs">
            {chartConfig[chart]?.label ?? chart}
        </span>
        <span class="text-lg font-bold leading-none sm:text-3xl">
            {total[key]?.toLocaleString?.() ?? 0}
        </span>
        </button>
    {/each}
</div>