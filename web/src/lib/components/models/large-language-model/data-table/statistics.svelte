<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	import type { LargeLanguageModel } from '../type';

	let { table }: { table: Table<LargeLanguageModel> } = $props();

	const filteredModels = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet Models()}
		{@const title = m.models()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:robot'}
		{@const models = filteredModels.length}
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
				<p class="text-7xl font-semibold">{models}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Models()}
</div>
