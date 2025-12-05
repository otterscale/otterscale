import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ColumnDef } from '@tanstack/table-core';

import type { Model } from '$lib/api/model/v1/model_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	modelName: m.model_name(),
	status: m.status(),
	description: m.description(),
	firstDeployedAt: m.first_deployed_at(),
	lastDeployedAt: m.last_deployed_at(),
	prefill: m.prefill(),
	decode: m.decode(),
	gpuRelation: m.gpu_relation(),
	test: m.test()
};

function getColumns(
	serviceUri: string,
	scope: string,
	reloadManager: ReloadManager
): ColumnDef<Model>[] {
	return [
		{
			id: 'expand',
			header: ({ table }) => {
				return renderSnippet(headers.row_expander, table);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.row_expander, row);
			},
			enableSorting: false,
			enableHiding: false
		},
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
			},
			enableHiding: false
		},
		{
			accessorKey: 'modelName',
			header: ({ column }) => {
				return renderSnippet(headers.modelName, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.modelName, row);
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
			accessorKey: 'status',
			header: ({ column }) => {
				return renderSnippet(headers.status, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.status, row);
			}
		},
		{
			accessorKey: 'description',
			header: ({ column }) => {
				return renderSnippet(headers.description, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.description, row);
			}
		},
		{
			accessorKey: 'prefill',
			header: ({ column }) => {
				return renderSnippet(headers.prefill, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.prefill, row);
			}
		},
		{
			accessorKey: 'decode',
			header: ({ column }) => {
				return renderSnippet(headers.decode, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.decode, row);
			}
		},
		{
			accessorKey: 'firstDeployedAt',
			header: ({ column }) => {
				return renderSnippet(headers.firstDeployedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.firstDeployedAt, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.firstDeployedAt,
					nextRow.original.firstDeployedAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		},
		{
			accessorKey: 'lastDeployedAt',
			header: ({ column }) => {
				return renderSnippet(headers.lastDeployedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.lastDeployedAt, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.lastDeployedAt,
					nextRow.original.lastDeployedAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		},
		{
			accessorKey: 'gpuRelation',
			header: ({ column }) => {
				return renderSnippet(headers.gpuRelation, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.gpuRelation, { row, scope, reloadManager });
			},
			enableHiding: false
		},
		{
			accessorKey: 'test',
			header: ({ column }) => {
				return renderSnippet(headers.test, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.test, { row, serviceUri: serviceUri });
			},
			enableHiding: false
		},
		{
			accessorKey: 'action',
			header: ({ column }) => {
				return renderSnippet(headers.action, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.action, { row, scope, reloadManager });
			},
			enableHiding: false
		}
	];
}
export { getColumns, messages };
