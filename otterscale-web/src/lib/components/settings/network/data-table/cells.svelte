<script lang="ts" module>
	import type { Network } from '$lib/api/network/v1/network_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import ViewVLAN from './view-vlan.svelte';
	import ViewSubnet from './view-subnet.svelte';
	import ViewIPAddresses from './view-ip-addresses.svelte';

	export const cells = {
		_row_picker: _row_picker,
		fabric,
		vlan,
		dhcpOn,
		subnet,
		ipAddresses,
		ipRanges,
		statistics
	};
</script>

{#snippet _row_picker(row: Row<Network>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet fabric(row: Row<Network>)}
	{#if row.original.fabric}
		{row.original.fabric.name}
	{/if}
{/snippet}

{#snippet vlan(row: Row<Network>)}
	{#if row.original.vlan}
		<span class="flex items-center">
			{row.original.vlan.name}
			<ViewVLAN vlan={row.original.vlan} />
		</span>
	{/if}
{/snippet}

{#snippet dhcpOn(row: Row<Network>)}
	{#if row.original.vlan}
		<Icon
			icon={row.original.vlan.dhcpOn ? 'ph:circle' : 'ph:x'}
			class={cn(row.original.vlan.dhcpOn ? 'text-primary' : 'text-destructive')}
		/>
	{/if}
{/snippet}

{#snippet subnet(row: Row<Network>)}
	{#if row.original.subnet}
		<span class="flex items-center">
			{row.original.subnet.name}
			<ViewSubnet subnet={row.original.subnet} />
		</span>
	{/if}
{/snippet}

{#snippet ipAddresses(row: Row<Network>)}
	<span class="flex justify-end">
		{#if row.original.subnet}
			<span class="flex items-center">
				{row.original.subnet.ipAddresses.length}
				<ViewIPAddresses ipAddresses={row.original.subnet.ipAddresses} />
			</span>
		{/if}
	</span>
{/snippet}

{#snippet ipRanges(row: Row<Network>)}
	<span class="flex justify-end">
		{#if row.original.subnet}
			{row.original.subnet.ipRanges.length}
		{/if}
	</span>
{/snippet}

{#snippet statistics(row: Row<Network>)}
	<span class="flex justify-end">
		{#if row.original.subnet && row.original.subnet.statistics}
			{row.original.subnet.statistics.usagePercent}
		{/if}
	</span>
{/snippet}
