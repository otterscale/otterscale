<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Connector } from './connector';
	import type { pbConnector } from '$lib/pb';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { stringify } from 'yaml';

	let {
		selected,
		connctors
	}: {
		selected: pbConnector | null;
		connctors: Connector[];
	} = $props();

	function getIcon(connectors: Connector[], key: string): string {
		return connectors.find((c) => c.key === key)?.icon ?? '';
	}
</script>

<div class="w-full flex-col items-center space-y-4">
	{#if selected}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class="flex w-full items-center justify-center hover:scale-105">
					<Icon icon={getIcon(connctors, selected.type)} class="size-12" />
					<div class="flex-col space-y-1 pl-4">
						<div class="text-sm text-foreground">{selected.name}</div>
						<div class="text-xs italic text-muted-foreground">{selected.id}</div>
					</div>
				</Tooltip.Trigger>
				<Tooltip.Content class="bg-slate-700 p-4">
					<span class="whitespace-break-spaces text-sm">{stringify(selected.workload.json)}</span>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
		<div class="grid grid-cols-3 items-center justify-center gap-2 px-2">
			<Badge variant="outline" class="col-span-1 w-min">Image</Badge>
			<p class="col-span-2">{selected.image}</p>
			<Badge variant="outline" class="col-span-1 w-min">Enabled</Badge>
			<p class="col-span-2">{selected.enabled}</p>
			<Badge variant="outline" class="col-span-1 w-min">Created</Badge>
			<p class="col-span-2">{selected.user}</p>
			{#if selected.workload}
				<Badge variant="outline" class="col-span-1 w-min">Updated</Badge>
				<p class="col-span-2">{selected.workload.user}</p>
				<Badge variant="outline" class="col-span-1 w-min">Lastest</Badge>
				<p class="col-span-2">{new Date(selected.workload.created).toLocaleString()}</p>
			{/if}
		</div>
	{:else}
		<div class="flex items-center justify-center space-x-4">
			<Skeleton class="size-12" />
			<div class="flex-col space-y-1">
				<Skeleton class="h-5 w-36" />
				<Skeleton class="h-5 w-36" />
			</div>
		</div>
		<div class="grid grid-cols-3 items-center justify-center gap-2 px-2">
			{#each Array(5) as _}
				<Skeleton class="col-span-1 h-5" />
				<Skeleton class="col-span-2 h-5" />
			{/each}
		</div>
	{/if}
</div>
