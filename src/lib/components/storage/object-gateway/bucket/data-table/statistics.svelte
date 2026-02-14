<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Card from '$lib/components/ui/card';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		table,
		serviceUri
	}: {
		table: Table<Bucket>;
		serviceUri: string;
	} = $props();

	const filteredBuckets = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
	const serviceIP = $derived(serviceUri.replace(/^https?:\/\//, ''));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet Buckets()}
		{@const title = m.bucket()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const buckets = filteredBuckets.length}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
					<Badge>
						{serviceIP}
					</Badge>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-7xl font-semibold">{buckets}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Buckets()}

	{#snippet Usage()}
		{@const title = m.usage()}
		{@const titleIcon = 'ph:chart-pie-bold'}
		{@const backgroundIcon = 'ph:disc'}
		{@const totalUsed = filteredBuckets
			.map((bucket) => Number(bucket.usedBytes))
			.reduce((a, current) => a + current, 0)}
		{@const { value: totalUsedValue, unit: totalUsedUnit } = formatCapacity(totalUsed)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-5xl font-semibold">{totalUsedValue} {totalUsedUnit}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Usage()}
</div>
