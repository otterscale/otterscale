<script lang="ts" module>
	import { type Writable } from 'svelte/store';

	import { type Release } from '$lib/api/application/v1/application_pb';
	import { type Chart } from '$lib/api/registry/v1/registry_pb';
	import { m } from '$lib/paraglide/messages';
	import { activeNamespace } from '$lib/stores';

	import ChartComponent from './chart.svelte';
	import FilterDeprecation from './filter-deprecation.svelte';
	import FilterKeyword from './filter-keyword.svelte';
	import FilterMaintainer from './filter-maintainer.svelte';
	import FilterName from './filter-name.svelte';
	import FilterReset from './filter-reset.svelte';
	import ImportChart from './import-chart.svelte';
	import Pagination from './pagination.svelte';
	import SynchronizeArtifactHub from './synchronize-artifact-hub.svelte';
	import Thumbnail from './thumbnail.svelte';
	import { FilterManager, PaginationManager } from './utils';
</script>

<script lang="ts">
	let {
		scope,
		charts,
		releases
	}: {
		scope: string;
		charts: Writable<Chart[]>;
		releases: Writable<Release[]>;
	} = $props();

	const releasesFromChartName = $derived(
		$releases
			.filter((release) => !$activeNamespace || release.namespace === $activeNamespace)
			.reduce((mapping, release) => {
				if (release.chart?.name) {
					if (!mapping.has(release.chart.name)) {
						mapping.set(release.chart.name, []);
					}
					mapping.get(release.chart.name)?.push(release);
				}
				return mapping;
			}, new Map<string, Release[]>())
	);
	const filterManager = $derived(new FilterManager($charts));
	const paginationManager = $derived(new PaginationManager(filterManager.filteredCharts));
</script>

<section class="mx-auto w-full space-y-4 bg-background">
	<div class="space-y-2 text-center">
		<h2 class="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
			{m.applications()}
		</h2>
		<p class="mx-auto text-lg text-muted-foreground">
			{m.store_description()}
		</p>
	</div>

	<div class="flex justify-between gap-4">
		<div class="flex items-center gap-1">
			<FilterName {filterManager} />
			<FilterKeyword {filterManager} />
			<FilterMaintainer {filterManager} />
			<FilterDeprecation {filterManager} />
			<FilterReset {filterManager} />
		</div>
		<div class="flex items-center gap-2">
			<ImportChart {scope} {charts} />
			<SynchronizeArtifactHub {scope} {charts} />
		</div>
	</div>

	<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		{#each filterManager.filteredCharts.slice(paginationManager.activePage * paginationManager.perPage, (paginationManager.activePage + 1) * paginationManager.perPage) as chart, index (index)}
			<ChartComponent
				{chart}
				chartReleases={releasesFromChartName.get(chart.name)}
				{scope}
				{releases}
			>
				<Thumbnail {chart} chartReleases={releasesFromChartName.get(chart.name)} />
			</ChartComponent>
		{/each}
	</div>

	<div class="text-center">
		<Pagination {paginationManager} />
	</div>
</section>
