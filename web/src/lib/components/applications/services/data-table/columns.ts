import type { ColumnDef } from '@tanstack/table-core';

import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import type { Service } from '../types';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	type: m.type(),
	clusterIp: m.cluster_ip(),
	ports: m.ports()
};

const columns: ColumnDef<Service>[] = [
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
		accessorKey: 'name',
		header: ({ column }) => {
			return renderSnippet(headers.name, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.name, row);
		}
	},
	{
		accessorKey: 'type',
		header: ({ column }) => {
			return renderSnippet(headers.type, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.type, row);
		}
	},
	{
		accessorKey: 'clusterIp',
		header: ({ column }) => {
			return renderSnippet(headers.clusterIp, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.clusterIp, row);
		}
	},
	{
		accessorKey: 'ports',
		header: ({ column }) => {
			return renderSnippet(headers.ports, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.ports, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.ports.length,
				nextRow.original.ports.length,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'endpoints',
		header: ({ column }) => {
			return renderSnippet(headers.endpoints, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.endpoints, row);
		},
		enableHiding: false
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
