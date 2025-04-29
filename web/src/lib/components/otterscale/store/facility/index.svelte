<script lang="ts">
	import { page } from '$app/state';
	import { cn } from '$lib/utils';

	// UI Components
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Skeleton } from '$lib/components/ui/skeleton';

	import { FacilityCreate, StoreFacility } from '$lib/components/otterscale/index';
	import { STORE_ITEMS_PER_ROW, STORE_ROWS_PER_PAGE } from '$lib/components/otterscale/index';
	import Filter from './filter.svelte';
	import Paging from './paging.svelte';
	import Search from './search.svelte';

	import type { Facility_Charm, Facility } from '$gen/api/nexus/v1/nexus_pb';

	let {
		charms,
		facilities
	}: {
		charms: Facility_Charm[];
		facilities: Facility[];
	} = $props();

	const ItemsPerPage = STORE_ITEMS_PER_ROW * STORE_ROWS_PER_PAGE;

	let charmToFacilities = $derived(
		facilities.reduce((m, f) => {
			const charmName = f.charmName;
			if (charmName) {
				if (!m.has(charmName)) {
					m.set(charmName, []);
				}
				m.get(charmName)?.push(f);
			}
			return m;
		}, new Map<string, Facility[]>())
	);

	let searchTerm = $state(page.url.searchParams.get('q') ?? '');
	let selectedCategories: string[] = $state([]);
	let onlyVerified = $state(false);
	let filteredCharms = $derived(
		charms.filter(
			(charm) =>
				charm.name.toLowerCase().includes(searchTerm.toLowerCase()) &&
				(selectedCategories.length === 0 ||
					selectedCategories.every((c) => charm.categories.includes(c))) &&
				(!onlyVerified || charm.verified)
		)
	);

	let activePage = $state(1);
</script>

<main class="grid justify-between">
	<div class="flex flex-col justify-between gap-2">
		<div class="w-full">
			<Search charms={filteredCharms} bind:searchTerm bind:activePage />
		</div>

		<div class="flex gap-2 overflow-auto">
			<span class="h-[660px] w-[13vw] min-w-[13vw] overflow-visible">
				<Filter {charms} bind:selectedCategories bind:onlyVerified bind:activePage />
			</span>
			<span
				class={cn(
					'grid h-[660px] w-full items-start justify-start gap-x-2',
					`grid-cols-${STORE_ITEMS_PER_ROW} grid-rows-${STORE_ROWS_PER_PAGE}`
				)}
			>
				{#each filteredCharms.slice((activePage - 1) * ItemsPerPage, activePage * ItemsPerPage) as filteredCharm}
					<AlertDialog.Root>
						<AlertDialog.Trigger>
							{@render CharmCard(filteredCharm)}
						</AlertDialog.Trigger>
						<AlertDialog.Content
							class={cn(
								'flex h-full flex-col justify-between',
								charmToFacilities.has(filteredCharm.name) ? 'min-w-[50vw]' : 'min-w-[38vw]'
							)}
						>
							<div class="flex h-full flex-col justify-between">
								<StoreFacility selectedCharm={filteredCharm} />
								<AlertDialog.Footer>
									<AlertDialog.Cancel class="mr-auto">Close</AlertDialog.Cancel>
									<AlertDialog.Action>
										<FacilityCreate />
									</AlertDialog.Action>
								</AlertDialog.Footer>
							</div>
						</AlertDialog.Content>
					</AlertDialog.Root>
				{/each}
			</span>
		</div>

		<span class="absolute -bottom-10 left-1/2 -translate-x-1/2">
			<Paging {filteredCharms} bind:activePage />
		</span>
	</div>
</main>

{#snippet CharmCard(filteredCharm: Facility_Charm)}
	<Card.Root
		class={cn(
			'duration-230 flex h-[320px] flex-col justify-between transition-transform hover:bg-muted'
		)}
	>
		<Card.Header class="h-[calc(320px*0.6*0.6*0.6)]">
			<Card.Title class="items-between flex gap-2">
				<Avatar.Root class="h-10 w-10">
					<Avatar.Image src={filteredCharm.icon} />
					<Avatar.Fallback>
						<Skeleton class="size-10" />
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="min-w-0 flex-1">
					<p class="truncate text-left text-base">{filteredCharm.name}</p>
					{#if filteredCharm.defaultArtifact}
						<span class="flex gap-1 overflow-x-auto text-muted-foreground">
							<p class="text-xs font-light">
								{filteredCharm.defaultArtifact?.channel}
							</p>
						</span>
					{/if}
				</span>
			</Card.Title>
			<Card.Description>
				{@render TruncatedBadges(filteredCharm.categories)}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="grid gap-2 p-4">
				<div class="text-left text-xs font-light">
					{@render TruncatedText(filteredCharm.description)}
				</div>
			</div>
		</Card.Content>
		<Card.Footer class="h-[calc(320px*0.6*0.6*0.6*0.6)]">
			{@render TruncatedBadges([filteredCharm.publisher])}
		</Card.Footer>
	</Card.Root>
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
