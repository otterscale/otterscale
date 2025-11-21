<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import { formatProgressColor } from '$lib/components/custom/progress/utils.svelte';
	import * as Card from '$lib/components/ui/card';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { formatCapacity, formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { table }: { table: Table<Image> } = $props();

	const filteredImages = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet Images()}
		{@const title = m.image()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const images = filteredImages.length}
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
				<p class="text-7xl font-semibold">{images}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Images()}

	{#snippet Usage()}
		{@const title = m.usage()}
		{@const titleIcon = 'ph:chart-pie-bold'}
		{@const backgroundIcon = 'ph:disc'}
		{@const totalQuota = filteredImages
			.map((image) => Number(image.quotaBytes))
			.reduce((a, current) => a + current, 0)}
		{@const totalUsed = filteredImages
			.map((image) => Number(image.usedBytes))
			.reduce((a, current) => a + current, 0)}
		{@const { value: totalQuotaValue, unit: totalQuotaUnit } = formatCapacity(totalQuota)}
		{@const { value: totalUsedValue, unit: totalUsedUnit } = formatCapacity(totalUsed)}
		{@const percentage = formatPercentage(totalUsed, totalQuota)}
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
						{totalUsedUnit} / {totalQuotaValue}
						{totalQuotaUnit}
					</p>
				</div>
			</Card.Content>
			<Progress
				value={Number(percentage ?? 0)}
				max={100}
				class={cn(
					formatProgressColor(totalUsed, totalQuota, 'STB'),
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
