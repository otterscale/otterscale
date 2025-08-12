<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/area/area-stock.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { fetchFlattenedRange } from '$lib/components/custom/prometheus';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Type definition for a single chart's data in the grid
	type ChartGridItem = { series: Map<string, SampleValue[] | undefined> };

	// Component's internal state
	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;
	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());
	let mounted = $state(false);
	let chartGridData: ChartGridItem[] = $state([]);
	let wsize: number = $derived(
		Math.ceil(chartGridData.length / 4) === 0 ? 2 : Math.ceil(chartGridData.length / 4)
	);

	// Constants for the query and layout
	const step = 1 * 60; // 1 minute step
	const MAX_CHARTS_IN_GRID = 24; // Maximum charts in the 4x6 grid
	const widthClasses: { [key: number]: string } = {
		2: 'w-1/2',
		3: 'w-1/3',
		4: 'w-1/4',
		5: 'w-1/5',
		6: 'w-1/6'
	};
	// The Prometheus query, derived from component props
	const query = $derived(
		`
		sum by (cpu) (rate(node_cpu_seconds_total{instance=~"${machine.fqdn}", mode!="idle"}[5m])) * 100
		`
	);

	/**
	 * A reactive effect that groups series data for each chart in the grid.
	 * If CPUs > 24, it wraps around and adds data to the first charts.
	 * This runs automatically whenever `serieses` changes.
	 */
	$effect(() => {
		const sortedKeys = Array.from(serieses.keys()).sort(
			(a, b) => parseInt(a.replace('cpu', ''), 10) - parseInt(b.replace('cpu', ''), 10)
		);

		if (sortedKeys.length === 0) {
			chartGridData = [];
			return;
		}

		// This will hold the series data for each chart in the grid.
		const newCharts: ChartGridItem[] = [];

		sortedKeys.forEach((key, index) => {
			const data = serieses.get(key);
			// Determine which chart this CPU belongs to using the modulo operator.
			const chartIndex = index % MAX_CHARTS_IN_GRID;

			if (index < MAX_CHARTS_IN_GRID) {
				// For the first 24 CPUs, create a new chart object in the array.
				newCharts[chartIndex] = { series: new Map([[key, data]]) };
			} else {
				// For CPUs 33 and beyond, add their data to the existing chart object.
				newCharts[chartIndex].series.set(key, data);
			}
		});

		chartGridData = newCharts;
	});
</script>

{#await fetchFlattenedRange(client, query, new Date(Date.now() - 60 * 60 * 1000), new Date(Date.now()), step)}
	<ComponentLoading />
{:then response}
	<Layout>
		{#snippet title()}
			<Title title="CPU" />
		{/snippet}

		{#snippet description()}
			<Description description="Core/Processor" />
		{/snippet}

		{#snippet content()}
			<Content data={response} timeRange={'1h'} />
		{/snippet}
	</Layout>
{:catch error}
	Error
{/await}
<!-- 
{#if mounted}
	<Template.Area title="CPU">
		{#snippet hint()}
			<p>Historical view of resource usage for all CPU core/Processor over the past hour</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Core/Processor</p>
		{/snippet}
		{#snippet content()}
			{#if chartGridData.length === 0}
				<Empty.Area />
			{:else}
				<div class="-mx-2 flex flex-wrap">
					{#each chartGridData as chart, i (i)}
						{@const currentSeriesMap = chart.series}
						{@const dataForChart = integrateSerieses(currentSeriesMap)}
						{@const seriesConfig = Array.from(currentSeriesMap.keys()).map((key) => {
							const cpuNumber = parseInt(key.replace('cpu', ''), 10) || 0;
							return {
								key: key,
								color: `hsl(${(Math.floor(cpuNumber / MAX_CHARTS_IN_GRID) * 60 + 203) % 360}, 83%, 60%)`
							};
						})}

						<div class="{widthClasses[wsize]} mb-2 px-2">
							<div
								class="mb-1 h-[60px] w-full resize border border-gray-100 last:mb-0 dark:border-gray-700"
							>
								<AreaChart
									data={dataForChart}
									x="time"
									yDomain={[0, 1]}
									series={seriesConfig}
									axis={false}
									grid={true}
									props={{
										tooltip: {
											root: {
												class:
													'bg-white/60 dark:bg-black/60 p-3 rounded-lg shadow-lg backdrop-blur-sm'
											},
											header: { class: 'font-light' },
											item: { format: 'percent' }
										}
									}}
									{renderContext}
									{debug}
								/>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		{/snippet}
	</Template.Area>
{:else}
	<ComponentLoading />
{/if} -->
