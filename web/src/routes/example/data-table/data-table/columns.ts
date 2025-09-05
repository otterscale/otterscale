import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { TableRow } from './type';

import { renderSnippet } from '$lib/components/ui/data-table/index.js';

const columns: ColumnDef<TableRow>[] = [
	{
		id: 'select',
		header: ({ table }) => {
			return renderSnippet(headers.row_picker, table);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.row_picker, row);
		},
		enableSorting: false,
		enableHiding: false,
	},
	{
		accessorKey: 'id',
		header: ({ column }) => {
			return renderSnippet(headers.id, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.id, row);
		},
		filterFn: (row, columnId, filterValues: Record<number, number> | undefined) => {
			if (!filterValues) {
				return true;
			}

			if (Object.values(filterValues).length != 2) {
				throw `invalid filter range for ${columnId}`;
			}

			const value = row.original.id;
			const range = Object.values(filterValues);
			range.sort();
			const [minimum, maximum] = range;

			return minimum <= value && value <= maximum;
		},
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderSnippet(headers.name, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.name, row);
		},
	},
	{
		accessorKey: 'isVerified',
		header: ({ column }) => {
			return renderSnippet(headers.isVerified, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.isVerified, row);
		},
	},
];

export { columns };
