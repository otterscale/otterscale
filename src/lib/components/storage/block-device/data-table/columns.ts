import type { ColumnDef } from '@tanstack/table-core';

import type { Image } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	poolName: m.pool_name(),
	usage: m.usage(),
	snapshots: m.snapshot()
};

function getColumns(scope: string, reloadManager: ReloadManager): ColumnDef<Image>[] {
	return [
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
			accessorKey: 'poolName',
			header: ({ column }) => {
				return renderSnippet(headers.poolName, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.poolName, row);
			},
			filterFn: 'arrIncludesSome'
		},
		{
			accessorKey: 'usage',
			header: ({ column }) => {
				return renderSnippet(headers.usage, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.usage, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					Number(previousRow.original.usedBytes) / Number(previousRow.original.quotaBytes),
					Number(nextRow.original.usedBytes) / Number(nextRow.original.quotaBytes),
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'snapshots',
			header: ({ column }) => {
				return renderSnippet(headers.snapshots, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.snapshots, { row, scope, reloadManager });
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.snapshots.length,
					nextRow.original.snapshots.length,
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
				return renderSnippet(cells.actions, { row, scope, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
