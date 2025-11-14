import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ColumnDef } from '@tanstack/table-core';

import type { Bucket } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	owner: m.owner(),
	usage: m.usage(),
	createTime: m.create_time()
};

function getColumns(scope: string, reloadManager: ReloadManager): ColumnDef<Bucket>[] {
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
			accessorKey: 'owner',
			header: ({ column }) => {
				return renderSnippet(headers.owner, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.owner, row);
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
					Number(previousRow.original.usedBytes),
					Number(nextRow.original.usedBytes),
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		},
		{
			accessorKey: 'createTime',
			header: ({ column }) => {
				return renderSnippet(headers.createTime, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.createTime, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.createdAt,
					nextRow.original.createdAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
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
