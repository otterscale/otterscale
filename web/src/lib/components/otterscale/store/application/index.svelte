<script lang="ts">
	import { page } from '$app/state';
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { StoreApplication, ReleaseCreate } from '$lib/components/otterscale/index';
	import { STORE_ITEMS_PER_ROW, STORE_ROWS_PER_PAGE } from '$lib/components/otterscale/index';
	import type { Application_Chart, Application_Release } from '$gen/api/nexus/v1/nexus_pb';
	import Paging from './paging.svelte';
	import Search from './search.svelte';
	import Filter from './filter.svelte';

	let {
		releases,
		charts
	}: {
		releases: Application_Release[];
		charts: Application_Chart[];
	} = $props();

	const ItemsPerPage = STORE_ITEMS_PER_ROW * STORE_ROWS_PER_PAGE;

	let releasesFromChart = $derived(
		releases.reduce((m, r) => {
			const chartName = r.chartName;
			if (chartName) {
				if (!m.has(chartName)) {
					m.set(chartName, []);
				}
				m.get(chartName)?.push(r);
			}
			return m;
		}, new Map<string, Application_Release[]>())
	);

	let searchTerm = $state(page.url.searchParams.get('q') ?? '');
	let filteredCharts = $derived(
		charts
			.filter((c) => !c.deprecated)
			.filter(
				(c) =>
					c.name.toLowerCase().includes(searchTerm.toLowerCase()) &&
					(selectedKeywords.length === 0 || selectedKeywords.every((k) => c.keywords.includes(k)))
			)
	);

	let selectedKeywords: string[] = $state([]);
	let activePage = $state(1);
</script>

<main class="grid justify-between">
	<div class="flex flex-col justify-between gap-2">
		<div class="w-full overflow-auto">
			{@render Keywords()}
			<div class="w-full">
				<Search {charts} bind:searchTerm bind:activePage />
			</div>
		</div>

		<div class="flex gap-2 overflow-auto">
			<span class="h-[660px] w-[13vw] min-w-[13vw] overflow-visible">
				<Filter charts={filteredCharts} bind:activePage bind:selectedKeywords />
			</span>
			<span
				class={cn(
					'grid h-[660px] w-full items-start justify-start gap-x-2',
					`grid-cols-${STORE_ITEMS_PER_ROW} grid-rows-${STORE_ROWS_PER_PAGE}`
				)}
			>
				{#each filteredCharts.slice((activePage - 1) * ItemsPerPage, activePage * ItemsPerPage) as filteredChart}
					<AlertDialog.Root>
						<AlertDialog.Trigger>
							{@render ChartCard(filteredChart)}
						</AlertDialog.Trigger>
						<AlertDialog.Content
							class={cn(
								'flex h-full flex-col justify-between',
								releasesFromChart.has(filteredChart.name) ? 'min-w-[62vw]' : 'min-w-[38vw]'
							)}
						>
							<div class="flex h-full flex-col justify-between">
								<StoreApplication
									selectedChart={filteredChart}
									selectedChartReleases={releasesFromChart.get(filteredChart.name)}
								/>
								<AlertDialog.Footer>
									<AlertDialog.Cancel class="mr-auto">Close</AlertDialog.Cancel>
									<AlertDialog.Action>
										<ReleaseCreate chart={filteredChart} />
									</AlertDialog.Action>
								</AlertDialog.Footer>
							</div>
						</AlertDialog.Content>
					</AlertDialog.Root>
				{/each}
			</span>
		</div>
		<span class="absolute -bottom-10 left-1/2 -translate-x-1/2">
			<Paging {filteredCharts} bind:activePage />
		</span>
	</div>
</main>

{#snippet ChartCard(filteredChart: Application_Chart)}
	<Card.Root
		class={cn(
			'duration-230 flex h-[320px] flex-col justify-between transition-transform hover:bg-muted'
		)}
	>
		<Card.Header class="h-[calc(320px*0.6*0.6*0.6)]">
			<Card.Title class="items-between flex gap-2">
				<Avatar.Root class="h-10 w-10">
					<Avatar.Image src={filteredChart.icon} />
					<Avatar.Fallback>
						<Skeleton class="size-10" />
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="min-w-0 flex-1">
					<p class="truncate text-left text-base">{filteredChart.name}</p>
					<span class="flex gap-1 overflow-x-auto text-muted-foreground">
						{#each filteredChart.versions as version}
							<p class="text-xs font-light">{version.applicationVersion}</p>
						{/each}
					</span>
				</span>
			</Card.Title>
			<Card.Description>
				{@render TruncatedBadges(filteredChart.keywords)}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<p class="p-4 text-left text-xs font-light">
				{@render TruncatedText(filteredChart.description)}
			</p>
		</Card.Content>
		<Card.Footer class="h-[calc(320px*0.6*0.6*0.6*0.6)]">
			{@render TruncatedBadges(filteredChart.maintainers.map((m) => m.name))}
		</Card.Footer>
	</Card.Root>
{/snippet}

{#snippet Keywords()}
	<span class="flex flex-wrap gap-2">
		{#each selectedKeywords as keyword}
			<Badge
				variant="secondary"
				class="cursor-pointer"
				onclick={() => (selectedKeywords = selectedKeywords.filter((k) => k !== keyword))}
			>
				{keyword}
				<Icon icon="ph:x" />
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet TruncatedText(text: string)}
	{@const Length = 100}
	{#if text.length > Length}
		{text.slice(0, 100 - 3)}...
	{:else}
		{text}
	{/if}
{/snippet}

{#snippet TruncatedBadges(badges: string[])}
	{@const Length = 1}
	<span class="flex gap-2 overflow-hidden">
		{#each badges.slice(0, Length) as keyword}
			<Badge variant="outline" class="text-xs">{keyword}</Badge>
		{/each}
		{#if badges.length > Length}
			<Badge variant="outline" class="whitespace-nowrap">+{badges.length - Length}</Badge>
		{/if}
	</span>
{/snippet}
