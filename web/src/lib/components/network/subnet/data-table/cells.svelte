<script lang="ts" module>
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { Network } from '$lib/api/network/v1/network_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { cn } from '$lib/utils';

	import ViewIPAddresses from './action-view-ip-addresses.svelte';
	import ViewSubnet from './action-view-subnet.svelte';
	import ViewVLAN from './action-view-vlan.svelte';
	import Actions from './cell-actions.svelte';
	import { ReservedIPRanges } from './cell-reserved-ip-ranges';

	export const cells = {
		row_picker,
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

{#snippet row_picker(row: Row<Network>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet fabric(row: Row<Network>)}
	<Layout.Cell class="items-start">
		{#if row.original.fabric}
			{row.original.fabric.name}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet vlan(row: Row<Network>)}
	<Layout.Cell class="items-start">
		{#if row.original.vlan}
			<div class="flex items-center gap-1">
				{row.original.vlan.name}
				<ViewVLAN vlan={row.original.vlan} />
			</div>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet dhcpOn(row: Row<Network>)}
	<Layout.Cell class="items-start">
		{#if row.original.vlan}
			<Icon
				icon={row.original.vlan.dhcpOn ? 'ph:circle' : 'ph:x'}
				class={cn(row.original.vlan.dhcpOn ? 'text-primary' : 'text-destructive')}
			/>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet subnet(row: Row<Network>)}
	<Layout.Cell class="items-start">
		{#if row.original.subnet}
			<span class="flex items-center">
				{row.original.subnet.name}
				<ViewSubnet subnet={row.original.subnet} />
			</span>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet ipAddresses(row: Row<Network>)}
	<Layout.Cell class="items-end">
		{#if row.original.subnet}
			<span class="flex items-center">
				{row.original.subnet.ipAddresses.length}
				<ViewIPAddresses ipAddresses={row.original.subnet.ipAddresses} />
			</span>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet ipRanges(data: { row: Row<Network>; reloadManager: ReloadManager })}
	<Layout.Cell class="items-end">
		{#if data.row.original.subnet}
			<ReservedIPRanges subnet={data.row.original.subnet} reloadManager={data.reloadManager} />
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet statistics(row: Row<Network>)}
	<Layout.Cell class="items-end">
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
	</Layout.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Network>; reloadManager: ReloadManager })}
	<Layout.Cell class="items-start">
		<Actions network={data.row.original} reloadManager={data.reloadManager} />
	</Layout.Cell>
{/snippet}
