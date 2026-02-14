<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Table } from '@tanstack/table-core';

	import type { Secret } from '$lib/api/application/v1/application_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Card from '$lib/components/ui/card';
	import ScrollArea from '$lib/components/ui/scroll-area/scroll-area.svelte';
	import { m } from '$lib/paraglide/messages';

	let { table }: { table: Table<Secret> } = $props();

	const filteredSecrets = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
	{#snippet Secrets()}
		{@const title = m.secrets()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:key'}
		{@const secrets = filteredSecrets.length}
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
				<p class="text-7xl font-semibold">{secrets}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Secrets()}

	{#snippet Types()}
		{@const title = m.type()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const types = new Set([...filteredSecrets.map((secret) => secret.type)])}
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
				<div class="flex gap-2">
					<p class="text-7xl font-semibold">{types.size}</p>
					<ScrollArea class="h-20 w-fit">
						<div class="space-y-1">
							{#each types as type (type)}
								<div class="h-5 p-1 text-xs">
									<Badge variant="outline">{type}</Badge>
								</div>
							{/each}
						</div>
					</ScrollArea>
				</div>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Types()}
</div>
