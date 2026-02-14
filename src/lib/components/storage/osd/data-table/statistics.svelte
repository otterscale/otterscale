<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import type { ObjectStorageDaemon } from '$lib/api/storage/v1/storage_pb';
	import { getProgressColor } from '$lib/components/custom/progress/utils.svelte';
	import * as Card from '$lib/components/ui/card';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { formatCapacity, formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { table }: { table: Table<ObjectStorageDaemon> } = $props();

	const filteredObjectStorageDaemons = $derived(
		table.getFilteredRowModel().rows.map((row) => row.original)
	);
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet ObjectStorageDaemons()}
		{@const title = m.osd()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const objectStorageDaemons = filteredObjectStorageDaemons.length}
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
				<p class="text-7xl font-semibold">{objectStorageDaemons}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render ObjectStorageDaemons()}

	{#snippet Usage()}
		{@const title = m.usage()}
		{@const titleIcon = 'ph:chart-pie-bold'}
		{@const backgroundIcon = 'ph:disc'}
		{@const totalSize = filteredObjectStorageDaemons
			.map((objectStorageDaemon) => Number(objectStorageDaemon.sizeBytes))
			.reduce((a, current) => a + current, 0)}
		{@const totalUsed = filteredObjectStorageDaemons
			.map((objectStorageDaemon) => Number(objectStorageDaemon.usedBytes))
			.reduce((a, current) => a + current, 0)}
		{@const { value: totalSizeValue, unit: totalSizeUnit } = formatCapacity(totalSize)}
		{@const { value: totalUsedValue, unit: totalUsedUnit } = formatCapacity(totalUsed)}
		{@const percentage = formatPercentage(totalUsed, totalSize, 0)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title>
					<div class="flex items-center gap-2 font-medium">
						<div
							class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
						>
							<Icon icon={titleIcon} class="size-5" />
						</div>
						<p class="font-bold">{title}</p>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<div class="space-y-1">
					<p class="text-5xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
					<p class="text-3xl text-muted-foreground">
						{totalUsedValue}
						{totalUsedUnit} / {totalSizeValue}
						{totalSizeUnit}
					</p>
				</div>
			</Card.Content>
			<Progress
				value={Number(percentage ?? 0)}
				max={100}
				class={cn(
					getProgressColor(totalUsed, totalSize, 'STB'),
					'absolute top-0 left-0 h-2 rounded-none'
				)}
			/>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Usage()}
</div>
