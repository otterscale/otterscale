import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ColumnDef } from '@tanstack/table-core';

import type { Secret } from '$lib/api/application/v1/application_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	type: m.type(),
	immutable: m.immutable(),
	labels: m.labels(),
	createdAt: m.create_time()
};

function getColumns(): ColumnDef<Secret>[] {
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
			accessorKey: 'namespace',
			header: ({ column }) => {
				return renderSnippet(headers.namespace, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.namespace, row);
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
			accessorKey: 'immutable',
			header: ({ column }) => {
				return renderSnippet(headers.immutable, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.immutable, row);
			}
		},
		{
			accessorKey: 'labels',
			header: ({ column }) => {
				return renderSnippet(headers.labels, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.labels, row);
			}
		},
		{
			accessorKey: 'createdAt',
			header: ({ column }) => {
				return renderSnippet(headers.created_at, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.created_at, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.createdAt,
					nextRow.original.createdAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		}
	];
}

export { getColumns, messages };
