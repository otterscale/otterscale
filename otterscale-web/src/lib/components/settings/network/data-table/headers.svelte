<script lang="ts" module>
	import type { Network } from '$lib/api/network/v1/network_pb';
	import { Layout, RowPickers, Sorter } from '$lib/components/custom/data-table';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
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

{#snippet _row_picker(table: Table<Network>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fabric(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>FABRIC</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet vlan(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>VLAN</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet dhcpOn(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>DHCP ON</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet subnet(column: Column<Network>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>SUBNET</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet ipAddresses(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>IP ADDRESS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ipRanges(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>IP RANGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet statistics(column: Column<Network>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>STATISTICS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<Network>)}{/snippet}
