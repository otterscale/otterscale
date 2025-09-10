import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	network: m.network(),
	node: m.node(),
	status: m.status(),
	cpu: m.cpu(),
	memory: m.memory(),
	disk: m.disk(),
};

const columns: ColumnDef<VirtualMachine>[] = [
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
		accessorKey: 'namespace',
		header: ({ column }) => {
			return renderSnippet(headers.namespace, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.namespace, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'network',
		header: ({ column }) => {
			return renderSnippet(headers.network, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.network, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'node',
		header: ({ column }) => {
			return renderSnippet(headers.node, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.node, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'status',
		header: ({ column }) => {
			return renderSnippet(headers.status, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.status, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'cpu',
		header: ({ column }) => {
			return renderSnippet(headers.cpu, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.cpu, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'memory',
		header: ({ column }) => {
			return renderSnippet(headers.memory, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.memory, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'disk',
		header: ({ column }) => {
			return renderSnippet(headers.disk, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.disk, row);
		},
		filterFn: 'arrIncludesSome',
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
