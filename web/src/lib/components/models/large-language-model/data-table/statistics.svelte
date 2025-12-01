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
		{@const modelsInTwoDays = filteredModels.filter((model) =>
			model.firstDeployedAt
				? timestampDate(model.firstDeployedAt) >= new Date(Date.now() - 2 * 24 * 60 * 60 * 1000)
				: false
		).length}
		{@const modelsInOneDay = filteredModels.filter((model) =>
			model.firstDeployedAt
				? timestampDate(model.firstDeployedAt) >= new Date(Date.now() - 1 * 24 * 60 * 60 * 1000)
				: false
		).length}
		{@const differenceInOneDay = modelsInOneDay + 1 - (modelsInTwoDays - modelsInOneDay)}

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
					{@const percentage = formatPercentage(differenceInOneDay, modelsInOneDay, 1)}
					{#if differenceInOneDay > 0}
						<Badge>
							<Icon icon="ph:trend-up" />
							<p>+{percentage}%</p>
						</Badge>
					{:else if differenceInOneDay < 0}
						<Badge>
							<Icon icon="ph:trend-down" />
							<p>-{percentage}%</p>
						</Badge>
					{:else}
						<Badge>
							<Icon icon="ph:minus" />
							<p>{percentage ?? 0}%</p>
						</Badge>
					{/if}
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex flex-col items-start gap-1.5 px-6 text-sm [.border-t]:pt-6">
				<div class="line-clamp-1 flex items-center gap-2 font-medium">
					{#if differenceInOneDay > 0}
						Trending up in the this day
						<Icon icon="ph:trend-up" />
					{:else if differenceInOneDay < 0}
						Trending down in the this day
						<Icon icon="ph:trend-down" />
					{:else}
						No change in the this day
					{/if}
				</div>
				<div class="text-muted-foreground">{modelsInOneDay} Models for the last 24 hours</div>
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
