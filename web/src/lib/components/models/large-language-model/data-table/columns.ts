import type { ColumnDef } from '@tanstack/table-core';

import { type LargeLangeageModel } from '../protobuf.svelte';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	version: m.version(),
	parameters: m.parameters(),
	accuracy: m.accuracy(),
	speed: m.speed(),
	architecture: m.architecture(),
	requests: m.requests(),
	uptime: m.uptime(),
};

const columns: ColumnDef<LargeLangeageModel>[] = [
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
		accessorKey: 'gpu_cache',
		header: ({ column }) => {
			return renderSnippet(headers.gpu_cache, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.gpu_cache, row);
		},
	},
	{
		accessorKey: 'kv_cache',
		header: ({ column }) => {
			return renderSnippet(headers.kv_cache, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.kv_cache, row);
		},
	},
	{
		accessorKey: 'requests',
		header: ({ column }) => {
			return renderSnippet(headers.requests, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.requests, row);
		},
	},
	{
		accessorKey: 'uptime',
		header: ({ column }) => {
			return renderSnippet(headers.uptime, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.uptime, row);
		},
	},
	{
		accessorKey: 'relation',
		header: ({ column }) => {
			return renderSnippet(headers.relation, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.relation, row);
		},
	},
];

export { columns, messages };
