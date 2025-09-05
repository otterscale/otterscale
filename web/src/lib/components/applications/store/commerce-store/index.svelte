<script lang="ts" module>
	import { type Writable } from 'svelte/store';

	import Chart from './chart.svelte';
	import FilterDeprecation from './filter-deprecation.svelte';
	import FilterKeyword from './filter-keyword.svelte';
	import FilterLicence from './filter-licence.svelte';
	import FilterMaintainer from './filter-maintainer.svelte';
	import FilterName from './filter-name.svelte';
	import FilterReset from './filter-reset.svelte';
	import Pagination from './pagination.svelte';
	import Thumbnail from './thumbnail.svelte';
	import Upload from './upload.svelte';
	import { FilterManager, PaginationManager } from './utils';

	import { type Application_Chart, type Application_Release } from '$lib/api/application/v1/application_pb';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		charts = $bindable(),
		releases = $bindable(),
	}: { charts: Writable<Application_Chart[]>; releases: Writable<Application_Release[]> } = $props();

	const releasesFromChartName = $derived(
		$releases.reduce((mapping, release) => {
			if (release.chartName) {
				if (!mapping.has(release.chartName)) {
					mapping.set(release.chartName, []);
				}
				mapping.get(release.chartName)?.push(release);
			}
			return mapping;
		}, new Map<string, Application_Release[]>()),
	);
	const filterManager = $derived(new FilterManager($charts));
	const paginationManager = $derived(new PaginationManager(filterManager.filteredCharts));
</script>

<section class="bg-background mx-auto w-full space-y-4">
	<div class="space-y-2 text-center">
		<h2 class="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
			{m.applications()}
		</h2>
		<p class="text-muted-foreground mx-auto text-lg">
			{m.store_description()}
		</p>
	</div>

	<div class="flex items-center gap-1">
		<FilterName {filterManager} />
		<FilterKeyword {filterManager} />
		<FilterMaintainer {filterManager} />
		<FilterLicence {filterManager} />
		<FilterDeprecation {filterManager} />
		<FilterReset {filterManager} />
		<Upload class="ml-auto" />
	</div>

	<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		{#each filterManager.filteredCharts.slice(paginationManager.activePage * paginationManager.perPage, (paginationManager.activePage + 1) * paginationManager.perPage) as chart}
			<Chart {chart} chartReleases={releasesFromChartName.get(chart.name)} bind:charts bind:releases>
				<Thumbnail {chart} chartReleases={releasesFromChartName.get(chart.name)} />
			</Chart>
		{/each}
	</div>

	<div class="text-center">
		<Pagination {paginationManager} />
	</div>
</section>
