import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	bus: m.bus(),
	bootOrder: m.boot_order(),
	dataVolume: m.data_volume(),
	type: m.type(),
	phase: m.phase(),
	boot: m.boot(),
	size: m.size(),
};

const columns: ColumnDef<EnhancedDisk>[] = [
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
		accessorKey: 'bus',
		header: ({ column }) => {
			return renderSnippet(headers.bus, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.bus, row);
		},
	},
	{
		accessorKey: 'bootOrder',
		header: ({ column }) => {
			return renderSnippet(headers.bootOrder, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.bootOrder, row);
		},
	},
	{
		accessorKey: 'dataVolume',
		header: ({ column }) => {
			return renderSnippet(headers.dataVolume, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.dataVolume, row);
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
		accessorKey: 'phase',
		header: ({ column }) => {
			return renderSnippet(headers.phase, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.phase, row);
		},
	},
	{
		accessorKey: 'boot',
		header: ({ column }) => {
			return renderSnippet(headers.boot, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.boot, row);
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
