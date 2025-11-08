import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { Network } from '$lib/api/network/v1/network_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	fabric: m.fabric(),
	vlan: m.vlan(),
	dhcpOn: m.dhcp_on(),
	subnet: m.subnet(),
	ipAddresses: m.ip_address(),
	ipRanges: m.ip_range(),
	statistics: m.statistics()
};

const columns: ColumnDef<Network>[] = [
	{
		id: 'select',
		header: ({ table }) => {
			return renderSnippet(headers.row_picker, table);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.row_picker, row);
		},
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'fabric',
		header: ({ column }) => {
			return renderSnippet(headers.fabric, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.fabric, row);
		},
		filterFn: (row, columnId, filterValue) => {
			if (!filterValue) {
				return true;
			}

			return row.original.fabric ? row.original.fabric.name.includes(filterValue) : false;
		}
	},
	{
		accessorKey: 'vlan',
		header: ({ column }) => {
			return renderSnippet(headers.vlan, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.vlan, row);
		},
		filterFn: (row, columnId, filterValue) => {
			if (!filterValue) {
				return true;
			}

			return row.original.vlan ? row.original.vlan.name.includes(filterValue) : false;
		}
	},
	{
		accessorKey: 'dhcpOn',
		header: ({ column }) => {
			return renderSnippet(headers.dhcpOn, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.dhcpOn, row);
		}
	},
	{
		accessorKey: 'subnet',
		header: ({ column }) => {
			return renderSnippet(headers.subnet, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.subnet, row);
		}
	},
	{
		accessorKey: 'ipAddresses',
		header: ({ column }) => {
			return renderSnippet(headers.ipAddresses, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.ipAddresses, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.subnet?.ipAddresses.length,
				nextRow.original.subnet?.ipAddresses.length,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'ipRanges',
		header: ({ column }) => {
			return renderSnippet(headers.ipRanges, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.ipRanges, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.subnet?.ipRanges?.length,
				nextRow.original.subnet?.ipRanges?.length,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'statistics',
		header: ({ column }) => {
			return renderSnippet(headers.statistics, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.statistics, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				Number(previousRow.original.subnet?.statistics?.available) /
					Number(previousRow.original.subnet?.statistics?.total),
				Number(nextRow.original.subnet?.statistics?.available) /
					Number(nextRow.original.subnet?.statistics?.total),
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'actions',
		header: ({ column }) => {
			return renderSnippet(headers.actions, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.actions, row);
		},
		enableHiding: false
	}
];

export { columns, messages };
