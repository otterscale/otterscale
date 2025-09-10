import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	description: m.description(),
	statusPhase: 'Status Phase',
	lastConditionMessage: 'Last Condition Message',
	lastConditionReason: 'Last Condition Reason',
};

const columns: ColumnDef<VirtualMachineSnapshot>[] = [
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
	},
	{
		accessorKey: 'statusPhase',
		header: ({ column }) => {
			return renderSnippet(headers.statusPhase, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.statusPhase, row);
		},
	},
	{
		accessorKey: 'description',
		header: ({ column }) => {
			return renderSnippet(headers.description, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.description, row);
		},
	},
	{
		accessorKey: 'lastConditionMessage',
		header: ({ column }) => {
			return renderSnippet(headers.lastConditionMessage, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.lastConditionMessage, row);
		},
	},
	{
		accessorKey: 'lastConditionReason',
		header: ({ column }) => {
			return renderSnippet(headers.lastConditionReason, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.lastConditionReason, row);
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
