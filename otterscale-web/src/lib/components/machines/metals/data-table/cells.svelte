<script lang="ts" module>
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { RowPickers } from '$lib/components/custom/data-table';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { dynamicPaths } from '$lib/path';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './actions.svelte';
	import Tags from './tags.svelte';

	export const cells = {
		_row_picker,
		fqdn_ip,
		powerState,
		status,
		cores_arch,
		ram,
		disk,
		storage,
		tags,
		actions
	};
</script>

{#snippet _row_picker(row: Row<Machine>)}
	<RowPickers.Cell {row} />
{/snippet}

{#snippet fqdn_ip(row: Row<Machine>)}
	<span class="flex items-center gap-1">
		{row.original.fqdn}

		<Button
			variant="ghost"
			size="icon"
			class="m-0 p-0"
			onclick={() =>
				goto(`${dynamicPaths.machinesMetal(page.params.scope).url}/${row.original.id}`)}
		>
			<Icon icon="ph:arrow-square-out" />
		</Button>
	</span>
	{#if row.original.ipAddresses}
		<span class="text-muted-foreground flex items-center gap-1 text-xs">
			{#each row.original.ipAddresses as ipAddress}
				{ipAddress}
			{/each}
		</span>
	{/if}
{/snippet}

{#snippet powerState(row: Row<Machine>)}
	<span class="flex items-center gap-1">
		<Icon
			icon={row.original.powerState === 'on' ? 'ph:power' : 'ph:power'}
			class={row.original.powerState === 'on' ? ' text-accent-foreground' : 'text-destructive'}
		/>
		<span class="flex flex-col items-start">
			{row.original.powerState}
			<p class="text-muted-foreground text-xs">{row.original.powerType}</p>
		</span>
	</span>
{/snippet}

{#snippet status(row: Row<Machine>)}
	{@const processingStates = [
		'commissioning',
		'deploying',
		'disk_erasing',
		'entering_rescue_mode',
		'exiting_rescue_mode',
		'releasing',
		'testing'
	]}
	<span class="flex flex-col">
		<Badge variant="outline">
			{row.original.status}
		</Badge>
		<span class="text-xs font-light">
			{#if row.original.statusMessage != 'Deployed'}
				<span class="flex items-center gap-1">
					{#if processingStates.includes(row.original.status.toLowerCase())}
						<Icon icon="ph:spinner" class="animate-spin" />
					{/if}
					<p class="invisible lg:visible">
						{row.original.statusMessage}
					</p>
				</span>
			{:else}
				<p class="invisible lg:visible">
					{`${row.original.osystem} ${row.original.hweKernel} ${row.original.distroSeries}`}
				</p>
			{/if}
		</span>
	</span>
{/snippet}

{#snippet cores_arch(row: Row<Machine>)}
	<span class="flex flex-col text-right">
		{row.original.cpuCount}
		<span class="text-muted-foreground">
			{row.original.architecture}
		</span>
	</span>
{/snippet}

{#snippet ram(row: Row<Machine>)}
	<span class="flex items-end justify-end gap-1">
		{formatCapacity(row.original.memoryMb).value}
		{formatCapacity(row.original.memoryMb).unit}
	</span>
{/snippet}

{#snippet disk(row: Row<Machine>)}
	<span class="flex items-end justify-end">
		{row.original.blockDevices.length}
	</span>
{/snippet}

{#snippet storage(row: Row<Machine>)}
	<span class="flex items-end justify-end gap-1">
		{formatCapacity(row.original.storageMb).value}
		{formatCapacity(row.original.storageMb).unit}
	</span>
{/snippet}

{#snippet tags(row: Row<Machine>)}
	<Tags machine={row.original} />
{/snippet}

{#snippet actions(row: Row<Machine>)}
	<Actions machine={row.original} />
{/snippet}
