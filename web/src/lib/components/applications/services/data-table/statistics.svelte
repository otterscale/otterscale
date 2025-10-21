<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Table } from '@tanstack/table-core';

	import { type Application_Service } from '$lib/api/application/v1/application_pb';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { currentKubernetes } from '$lib/stores';

	let { table }: { table: Table<Application_Service> } = $props();
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
	<Card.Root class="relative overflow-hidden">
		<Card.Header class="gap-3">
			<Card.Title class="flex items-center gap-2 font-medium">
				<div class="bg-primary/10 text-primary flex size-8 shrink-0 items-center justify-center rounded-md">
					<Icon icon="ph:chart-bar" />
				</div>
				Services
			</Card.Title>
			<Card.Description>
				<Badge>
					over scope {$currentKubernetes?.scope}
				</Badge>
			</Card.Description>
		</Card.Header>
		<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
			<p class="text-5xl font-semibold">{table.getCoreRowModel().rows.length}</p>
		</Card.Content>
		<div
			class="text-primary/5 absolute top-0 -right-16 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
		>
			<Icon icon="ph:squares-four" class="size-72" />
		</div>
	</Card.Root>
	<Card.Root class="relative overflow-hidden">
		<Card.Header class="gap-3">
			<Card.Title class="flex items-center gap-2 font-medium">
				<div class="bg-primary/10 text-primary flex size-8 shrink-0 items-center justify-center rounded-md">
					<Icon icon="ph:chart-pie-slice" />
				</div>
				Types
			</Card.Title>
			<Card.Description>
				<Badge>
					over scope {$currentKubernetes?.scope}
				</Badge>
			</Card.Description>
		</Card.Header>
		<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
			<p class="text-5xl font-semibold">
				{new Set([...table.getCoreRowModel().rows.map((row) => row.getValue('type'))]).size}
			</p>
		</Card.Content>
		<div
			class="text-primary/5 absolute top-0 -right-16 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
		>
			<Icon icon="ph:broadcast" class="size-72" />
		</div>
	</Card.Root>
</div>
