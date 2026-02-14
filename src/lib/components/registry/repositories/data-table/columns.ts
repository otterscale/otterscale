import type { ColumnDef } from '@tanstack/table-core';

import type { Repository } from '$lib/api/registry/v1/registry_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	manifests: m.manifests(),
	sizeBytes: m.size(),
	latestTag: m.latest_tag()
};

function getColumns(scope: string, reloadManager: ReloadManager): ColumnDef<Repository>[] {
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
			accessorKey: 'latestTag',
			header: ({ column }) => {
				return renderSnippet(headers.latestTag, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.latestTag, row);
			}
		},
		{
			accessorKey: 'sizeBytes',
			header: ({ column }) => {
				return renderSnippet(headers.sizeBytes, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.sizeBytes, row);
			}
		},

		{
			accessorKey: 'manifests',
			header: ({ column }) => {
				return renderSnippet(headers.manifests, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.manifests, { row, scope, reloadManager });
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.manifestCount,
					nextRow.original.manifestCount,
					(p, n) => p < n,
					(p, n) => p === n
				)
		}
	];
}

export { getColumns, messages };
