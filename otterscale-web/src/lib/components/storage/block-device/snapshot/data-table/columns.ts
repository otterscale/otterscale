import type { Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import type { ColumnDef } from '@tanstack/table-core';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Image_Snapshot>[] = [
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
		accessorKey: 'protect',
		header: ({ column }) => {
			return renderSnippet(headers.protect, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.protect, row);
		},
		sortingFn: (previousRow, nextRow, columnId) =>
			getSortingFunction(
				previousRow.original.protected,
				nextRow.original.protected,
				(p, n) => Number(p) < Number(n),
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'usage',
		header: ({ column }) => {
			return renderSnippet(headers.usage, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.usage, row);
		},
		sortingFn: (previousRow, nextRow, columnId) =>
			getSortingFunction(
				Number(previousRow.original.usedBytes) / Number(previousRow.original.quotaBytes),
				Number(nextRow.original.usedBytes) / Number(nextRow.original.quotaBytes),
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
		}
	}
];

export { columns };
