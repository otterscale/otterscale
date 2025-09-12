import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { VirtualMachineDisk } from '$lib/api/kubevirt/v1/kubevirt_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	type: m.type(),
	bus: 'bus',
	source: m.source(),
	sourceType: 'source type',
	size: 'disk size',
};

const columns: ColumnDef<VirtualMachineDisk>[] = [
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
		accessorKey: 'type',
		header: ({ column }) => {
			return renderSnippet(headers.type, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.type, row);
		},
	},
	{
		accessorKey: 'bus',
		header: ({ column }) => {
			return renderSnippet(headers.bus, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.bus, row);
		},
	},
	{
		accessorKey: 'sourceType',
		header: ({ column }) => {
			return renderSnippet(headers.sourceType, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.sourceType, row);
		},
	},
	{
		accessorKey: 'source',
		header: ({ column }) => {
			return renderSnippet(headers.source, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.source, row);
		},
	},
	{
		accessorKey: 'size',
		header: ({ column }) => {
			return renderSnippet(headers.size, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.size, row);
		},
	},
	// {
	// 	accessorKey: 'actions',
	// 	header: ({ column }) => {
	// 		return renderSnippet(headers.actions, column);
	// 	},
	// 	cell: ({ row }) => {
	// 		return renderSnippet(cells.actions, row);
	// 	},
	// 	enableHiding: false,
	// },
];

export { columns, messages };
