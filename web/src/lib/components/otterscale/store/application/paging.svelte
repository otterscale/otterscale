<script lang="ts">
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { STORE_ITEMS_PER_ROW, STORE_ROWS_PER_PAGE } from '$lib/components/otterscale/index';
	import { type Application_Chart } from '$gen/api/application/v1/application_pb';

	let {
		filteredCharts = $bindable(),
		activePage = $bindable()
	}: {
		filteredCharts: Application_Chart[];
		activePage: number;
	} = $props();

	const ItemsPerPage = STORE_ITEMS_PER_ROW * STORE_ROWS_PER_PAGE;
</script>

<Pagination.Root count={filteredCharts.length} perPage={ItemsPerPage} bind:page={activePage}>
	{#snippet children({ pages })}
		<Pagination.Content class="rounded-lg bg-muted p-1">
			<Pagination.Item>
				<Pagination.PrevButton />
			</Pagination.Item>
			{#each pages as page (page.key)}
				{#if page.type === 'ellipsis'}
					<Pagination.Item>
						<Pagination.Ellipsis />
					</Pagination.Item>
				{:else}
					<Pagination.Item>
						<Pagination.Link {page} isActive={activePage === page.value}>
							{page.value}
						</Pagination.Link>
					</Pagination.Item>
				{/if}
			{/each}
			<Pagination.Item>
				<Pagination.NextButton />
			</Pagination.Item>
		</Pagination.Content>
	{/snippet}
</Pagination.Root>
