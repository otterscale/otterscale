<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let { table }: { table: Table<Application_Service_Port> } = $props();

	const filteredPorts = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
	{#snippet Ports()}
		{@const title = m.port()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:network'}
		{@const ports = filteredPorts.length}
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
				<p class="text-7xl font-semibold">{ports}</p>
			</Card.Content>
			<div
				class="absolute top-8 -right-8 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-48" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Ports()}

	{#snippet NodePorts()}
		{@const title = m.node_port()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:share-network'}
		{@const nodePorts = filteredPorts.filter((port) => port.nodePort).length}
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
				<p class="text-7xl font-semibold">{nodePorts}</p>
			</Card.Content>
			<div
				class="absolute top-8 -right-8 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-48" />
			</div>
		</Card.Root>
	{/snippet}
	{@render NodePorts()}
</div>
