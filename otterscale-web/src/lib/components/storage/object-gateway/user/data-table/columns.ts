import type { User } from '$lib/api/storage/v1/storage_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import type { ColumnDef } from '@tanstack/table-core';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<User>[] = [
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
		accessorKey: 'id',
		header: ({ column }) => {
			return renderSnippet(headers.id, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.id, row);
		}
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderSnippet(headers.name, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.name, row);
		},
		filterFn: 'arrIncludesSome'
	},
	{
		accessorKey: 'suspended',
		header: ({ column }) => {
			return renderSnippet(headers.suspended, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.suspended, row);
		},
		filterFn: 'equals'
	},
	{
		accessorKey: 'keys',
		header: ({ column }) => {
			return renderSnippet(headers.keys, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.keys, row);
		}
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

export { columns };
