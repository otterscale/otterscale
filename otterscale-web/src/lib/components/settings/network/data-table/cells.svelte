<script lang="ts" module>
	import type { Network } from '$lib/api/network/v1/network_pb';
	import { RowPickers } from '$lib/components/custom/data-table';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './actions.svelte';
	import { ReservedIPRanges } from './reservedIPRanges';
	import ViewIPAddresses from './view-ip-addresses.svelte';
	import ViewSubnet from './view-subnet.svelte';
	import ViewVLAN from './view-vlan.svelte';

	export const cells = {
		_row_picker,
		fabric,
		vlan,
		dhcpOn,
		subnet,
		ipAddresses,
		ipRanges,
		statistics,
		actions
	};
</script>

{#snippet _row_picker(row: Row<Network>)}
	<RowPickers.Cell {row} />
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
			<ReservedIPRanges subnet={row.original.subnet} />
		{/if}
	</span>
{/snippet}

{#snippet statistics(row: Row<Network>)}
	{#if row.original.subnet && row.original.subnet.statistics}
		<Progress.Root
			numerator={Number(row.original.subnet.statistics.available)}
			denominator={Number(row.original.subnet.statistics.total)}
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{numerator}/{denominator}
			{/snippet}
		</Progress.Root>
	{/if}
{/snippet}

{#snippet actions(row: Row<Network>)}
	<Actions network={row.original}></Actions>
{/snippet}
