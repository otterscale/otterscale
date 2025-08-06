<script lang="ts" module>
	import {
		ApplicationService,
		type Application_Chart
	} from '$lib/api/application/v1/application_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Chart } from './chart';
	import FilterKeyword from './filter-keyword.svelte';
	import FilterMaintainer from './filter-maintainer.svelte';
	import FilterName from './filter-name.svelte';
	import FilterReset from './filter-reset.svelte';
	import Pagination from './pagination.svelte';
	import Thumbnail from './thumbnail.svelte';
	import Upload from './upload.svelte';
	import { FilterManager, PaginationManager } from './utils.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const charts = writable<Application_Chart[]>([]);

	let isMounted = $state(false);
	onMount(async () => {
		await applicationClient
			.listCharts({})
			.then((response) => {
				charts.set(response.charts);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
	});

	const filterManager = $derived(new FilterManager($charts));
	const paginationManager = $derived(new PaginationManager(filterManager.filteredCharts));
</script>

<section class="bg-background container mx-auto space-y-6 py-2 md:py-4">
	<div class="space-y-2 text-center">
		<h2 class="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">Applications</h2>
		<p class="text-muted-foreground mx-auto text-lg">
			Browse and install verified applications for your cluster
		</p>
	</div>

	<div class="flex items-center gap-1">
		<FilterName {filterManager} />
		<FilterKeyword {filterManager} />
		<FilterMaintainer {filterManager} />
		<FilterReset {filterManager} />
		<Upload class="ml-auto" />
	</div>

	<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		{#each filterManager.filteredCharts.slice(paginationManager.activePage * paginationManager.perPage, (paginationManager.activePage + 1) * paginationManager.perPage) as chart}
			<Chart {chart}>
				<Thumbnail {chart} />
			</Chart>
		{/each}
	</div>

	<div class="text-center">
		<Pagination {paginationManager} />
	</div>
</section>
