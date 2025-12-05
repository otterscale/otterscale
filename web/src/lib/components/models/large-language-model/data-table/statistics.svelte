<script lang="ts">
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Card from '$lib/components/ui/card';
	import { formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { table }: { table: Table<Model> } = $props();

	const filteredModels = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet Models()}
		{@const title = m.models()}
		{@const backgroundIcon = 'ph:robot'}
		{@const models = filteredModels.length}
		{@const pastModels = filteredModels.filter((model) =>
			model.firstDeployedAt
				? timestampDate(model.firstDeployedAt) >= new Date(Date.now() - 14 * 24 * 60 * 60 * 1000) &&
					timestampDate(model.firstDeployedAt) < new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
				: false
		).length}
		{@const recentModels = filteredModels.filter((model) =>
			model.firstDeployedAt
				? timestampDate(model.firstDeployedAt) >= new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
				: false
		).length}
		{@const difference = recentModels - pastModels}
		<Card.Root
			class="@container/card relative flex flex-col gap-6 overflow-hidden rounded-xl border bg-card py-6 text-card-foreground shadow-sm"
		>
			<Card.Header class="gap-3">
				<Card.Description class="text-sm text-muted-foreground">
					{title}
				</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					{models}
				</Card.Title>
				<Card.Action class="col-start-2 row-span-2 row-start-1 self-start justify-self-end">
					{@const percentage = formatPercentage(difference, pastModels, 1)}
					{#if percentage}
						<Badge>
							{#if difference > 0}
								<Icon icon="ph:trend-up" />
							{:else if difference < 0}
								<Icon icon="ph:trend-down" />
							{:else}
								<Icon icon="ph:minus" />
							{/if}
							<p>{percentage}%</p>
						</Badge>
					{/if}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex flex-col items-start gap-1.5 px-6 text-sm [.border-t]:pt-6">
				<div class="line-clamp-1 flex items-center gap-2 font-medium">
					{#if difference > 0}
						Trending up in 7 day
						<Icon icon="ph:trend-up" />
					{:else if difference < 0}
						Trending down in 7 day
						<Icon icon="ph:trend-down" />
					{:else}
						No change in 7 day
					{/if}
				</div>
				<div class="text-muted-foreground">{recentModels} Models for the last 14 days</div>
			</Card.Footer>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Models()}
</div>
