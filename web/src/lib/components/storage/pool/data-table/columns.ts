import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { Pool } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	poolType: m.type(),
	applications: m.applications(),
	placementGroupState: m.placement_group_state(),
	usage: m.usage(),
	iops: m.iops(),
};

const columns: ColumnDef<Pool>[] = [
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
		accessorKey: 'name',
		header: ({ column }) => {
			return renderSnippet(headers.name, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.name, row);
		},
	},
	{
		accessorKey: 'poolType',
		header: ({ column }) => {
			return renderSnippet(headers.type, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.type, row);
		},
	},
	{
		accessorKey: 'applications',
		header: ({ column }) => {
			return renderSnippet(headers.applications, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.applications, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'placementGroupState',
		header: ({ column }) => {
			return renderSnippet(headers.placement_group_state, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.placement_group_state, row);
		},
		filterFn: (row, columnId, filterValue) => {
			const values = Object.keys(row.getValue(columnId) ?? {});
			if (!values.length || !filterValue.length) return true;
			return values.some((value) => filterValue.includes(value));
		},
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
				Number(previousRow.original.quotaBytes) !== 0
					? Number(previousRow.original.usedBytes) / Number(previousRow.original.quotaBytes)
					: 0,
				Number(nextRow.original.quotaBytes) !== 0
					? Number(nextRow.original.usedBytes) / Number(nextRow.original.quotaBytes)
					: 0,
				(p, n) => p < n,
				(p, n) => p === n,
			),
	},
	{
		accessorKey: 'iops',
		header: ({ column }) => {
			return renderSnippet(headers.iops, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.iops, row);
		},
	},
	{
		accessorKey: 'actions',
		header: ({ column }) => {
			return renderSnippet(headers.actions, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.actions, row);
		},
		enableHiding: false,
	},
];

export { columns, messages };
