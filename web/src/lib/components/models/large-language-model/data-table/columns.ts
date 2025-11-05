import type { ColumnDef } from '@tanstack/table-core';

import { type LargeLanguageModel, type Meta } from '../type';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	chart_version: m.version(),
	application_version: m.version(),
	status: m.status(),
	description: m.description(),
	requests: m.requests(),
	limits: m.limits(),
	first_deployed_at: m.time(),
	last_deployed_at: m.time(),
	relation: m.relational(),
	pods: m.pods(),
	kv_cache: m.kv_cache(),
};

const columns: ColumnDef<LargeLanguageModel>[] = [
	{
		id: 'expand',
		header: ({ table }) => renderSnippet(headers.row_expander, table),
		cell: ({ row }) => renderSnippet(cells.row_expander, row),
		enableSorting: false,
		enableHiding: false,
		meta: {
			isRowAction: true,
		} as Meta,
	},
	{
		id: 'select',
		header: ({ table }) => renderSnippet(headers.row_picker, table),
		cell: ({ row }) => renderSnippet(cells.row_picker, row),
		enableSorting: false,
		enableHiding: false,
		meta: {
			isRowAction: true,
		} as Meta,
	},
	{
		accessorKey: 'name',
		header: ({ column }) => renderSnippet(headers.name, column),
		cell: ({ row }) => renderSnippet(cells.name, row),
	},
	{
		accessorKey: 'namespace',
		header: ({ column }) => renderSnippet(headers.namespace, column),
		cell: ({ row }) => renderSnippet(cells.namespace, row),
	},
	{
		accessorKey: 'chart_version',
		header: ({ column }) => renderSnippet(headers.chart_version, column),
		cell: ({ row }) => renderSnippet(cells.chart_version, row),
	},
	{
		accessorKey: 'application_version',
		header: ({ column }) => renderSnippet(headers.application_version, column),
		cell: ({ row }) => renderSnippet(cells.application_version, row),
	},
	{
		accessorKey: 'status',
		header: ({ column }) => renderSnippet(headers.status, column),
		cell: ({ row }) => renderSnippet(cells.status, row),
	},
	{
		accessorKey: 'description',
		header: ({ column }) => renderSnippet(headers.description, column),
		cell: ({ row }) => renderSnippet(cells.description, row),
	},
	{
		accessorKey: 'requests',
		header: ({ column }) => renderSnippet(headers.requests, column),
		cell: ({ row }) => renderSnippet(cells.requests, row),
	},
	{
		accessorKey: 'limits',
		header: ({ column }) => renderSnippet(headers.limits, column),
		cell: ({ row }) => renderSnippet(cells.limits, row),
	},
	{
		accessorKey: 'first_deployed_at',
		header: ({ column }) => renderSnippet(headers.first_deployed_at, column),

		cell: ({ row }) => renderSnippet(cells.first_deployed_at, row),
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.firstDeployedAt,
				nextRow.original.firstDeployedAt,
				(p, n) => Date.parse(String(p)) < Date.parse(String(n)),
				(p, n) => Date.parse(String(p)) === Date.parse(String(n)),
			),
	},
	{
		accessorKey: 'last_deployed_at',
		header: ({ column }) => renderSnippet(headers.last_deployed_at, column),
		cell: ({ row }) => renderSnippet(cells.last_deployed_at, row),
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.lastDeployedAt,
				nextRow.original.lastDeployedAt,
				(p, n) => Date.parse(String(p)) < Date.parse(String(n)),
				(p, n) => Date.parse(String(p)) === Date.parse(String(n)),
			),
	},
	{
		accessorKey: 'relation',
		header: ({ column }) => renderSnippet(headers.relation, column),
		cell: ({ row }) => renderSnippet(cells.relation, row),
		enableSorting: false,
	},
	{
		accessorKey: 'pods',
		header: ({ column }) => renderSnippet(headers.pods, column),
		cell: ({ row }) => renderSnippet(cells.pods, row),
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.pods.length,
				nextRow.original.pods.length,
				(p, n) => p < n,
				(p, n) => p === n,
			),
	},
	{
		accessorKey: 'action',
		header: ({ column }) => renderSnippet(headers.action, column),
		cell: ({ row }) => renderSnippet(cells.action, row),
		enableHiding: false,
	},
];

export { columns, messages };
