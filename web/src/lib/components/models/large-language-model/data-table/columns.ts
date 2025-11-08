import type { ColumnDef } from '@tanstack/table-core';

import { type LargeLanguageModel } from '../type';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	gpu_cache: m.gpu_cache(),
	application: m.application(),
	replicas: m.replica(),
	healthies: m.health(),
	kv_cache: m.kv_cache(),
	requests: m.requests(),
	time_to_first_token: m.uptime()
};

const columns: ColumnDef<LargeLanguageModel>[] = [
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
		accessorKey: 'model',
		header: ({ column }) => {
			return renderSnippet(headers.model, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.model, row);
		}
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
		accessorKey: 'replicas',
		header: ({ column }) => {
			return renderSnippet(headers.replicas, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.replicas, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.application.replicas,
				nextRow.original.application.replicas,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'healthies',
		header: ({ column }) => {
			return renderSnippet(headers.healthies, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.healthies, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.application.healthies,
				nextRow.original.application.healthies,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'gpu_cache',
		header: ({ column }) => {
			return renderSnippet(headers.gpu_cache, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.gpu_cache, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.metrics.gpu_cache,
				nextRow.original.metrics.gpu_cache,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'kv_cache',
		header: ({ column }) => {
			return renderSnippet(headers.kv_cache, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.kv_cache, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.metrics.kv_cache,
				nextRow.original.metrics.kv_cache,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'requests',
		header: ({ column }) => {
			return renderSnippet(headers.requests, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.requests, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.metrics.requests,
				nextRow.original.metrics.requests,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'time_to_first_token',
		header: ({ column }) => {
			return renderSnippet(headers.time_to_first_token, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.time_to_first_token, row);
		},
		sortingFn: (previousRow, nextRow) =>
			getSortingFunction(
				previousRow.original.metrics.time_to_first_token,
				nextRow.original.metrics.time_to_first_token,
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'relation',
		header: ({ column }) => {
			return renderSnippet(headers.relation, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.relation, row);
		},
		enableHiding: false
	},

	{
		accessorKey: 'action',
		header: ({ column }) => {
			return renderSnippet(headers.action, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.action, row);
		},
		enableHiding: false
	}
];

export { columns, messages };
