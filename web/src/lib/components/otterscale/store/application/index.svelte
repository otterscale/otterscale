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
			<span class="overflow-visible">
				<Filter charts={filteredCharts} bind:activePage bind:selectedKeywords />
			</span>
			<span
				class={cn(
					'grid w-full items-start justify-start gap-x-2',
					`grid-cols-${STORE_ITEMS_PER_ROW} grid-rows-${STORE_ROWS_PER_PAGE}`
				)}
			>
				{#each filteredCharts.slice((activePage - 1) * ItemsPerPage, activePage * ItemsPerPage) as filteredChart}
					<AlertDialog.Root>
						<AlertDialog.Trigger>
							{@render ChartCard(filteredChart)}
						</AlertDialog.Trigger>
						<AlertDialog.Content
							interactOutsideBehavior="close"
							class={cn(
								'flex flex-col justify-between',
								releasesFromChart.has(filteredChart.name) ? 'min-w-[62vw]' : 'min-w-[38vw]'
							)}
						>
							<div class="flex h-full flex-col justify-between">
								<StoreApplication
									bind:releases
									selectedChart={filteredChart}
									selectedChartReleases={releasesFromChart.get(filteredChart.name)}
								/>
								<AlertDialog.Footer>
									<AlertDialog.Cancel class="mr-auto">Close</AlertDialog.Cancel>
									<AlertDialog.Action>
										<ReleaseCreate bind:releases chart={filteredChart} />
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
	<Card.Root class={cn('flex-col transition-transform hover:bg-muted')}>
		<Card.Header>
			<Card.Title class="flex items-center gap-2 space-x-2">
				<Avatar.Root class="h-12 w-12">
					<Avatar.Image src={filteredChart.icon} />
					<Avatar.Fallback>
						<Skeleton class="size-12" />
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="flex-col space-y-1">
					<p class="truncate text-left text-base">
						{filteredChart.name
							.split('-')
							.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
							.join(' ')}
					</p>
					<span class="flex text-muted-foreground">
						{#each filteredChart.versions as version}
							<p class="text-xs font-light">{version.applicationVersion}</p>
						{/each}
					</span>
				</span>
			</Card.Title>
		</Card.Header>
		<Card.Content>
			<p class="text-left text-sm">
				{@render TruncatedText(filteredChart.description)}
			</p>
		</Card.Content>
		<Card.Footer class="space-x-2">
			{#if filteredChart.deprecated}
				<Badge variant="destructive">deprecated</Badge>
			{/if}
			{#if releasesFromChart.has(filteredChart.name)}
				<Badge>Installed</Badge>
			{/if}
			{@render TruncatedBadges(filteredChart.keywords)}
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
	<span class="flex space-x-2 overflow-hidden">
		{#each badges.slice(0, Length) as keyword}
			<Badge variant="outline">{keyword}</Badge>
		{/each}
		{#if badges.length > Length}
			<Badge variant="outline" class="whitespace-nowrap">+{badges.length - Length}</Badge>
		{/if}
	</span>
{/snippet}
