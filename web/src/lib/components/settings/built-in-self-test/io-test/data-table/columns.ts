import { type ColumnDef } from '@tanstack/table-core';

import { type TestResult } from '$lib/api/configuration/v1/configuration_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { ReloadManager } from '$lib/components/custom/reloader';

const messages = {
	name: m.name(),
	status: m.status(),
	target: m.target(),
	accessMode: m.access_mode(),
	jobCount: m.job_count(),
	runTime: m.run_time(),
	blockSize: m.block_size(),
	fileSize: m.file_size(),
	ioDepth: m.io_depth(),
	bandwidth: m.bandwidth(),
	iops: m.iops(),
	latencyMinimum: m.latency_minimum(),
	latencyMaximum: m.latency_maximum(),
	latencyMean: m.latency_mean(),
	createdBy: m.created_by(),
	startedAt: m.started_at(),
	completedAt: m.completed_at()
};

function getColumns(reloadManager: ReloadManager): ColumnDef<TestResult>[] {
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
			accessorKey: 'status',
			header: ({ column }) => {
				return renderSnippet(headers.status, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.status, row);
			}
		},
		{
			accessorKey: 'target',
			header: ({ column }) => {
				return renderSnippet(headers.target, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.target, row);
			}
		},
		{
			accessorKey: 'accessMode',
			header: ({ column }) => {
				return renderSnippet(headers.accessMode, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.accessMode, row);
			}
		},
		{
			accessorKey: 'jobCount',
			header: ({ column }) => {
				return renderSnippet(headers.jobCount, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.jobCount, row);
			}
		},
		{
			accessorKey: 'runTime',
			header: ({ column }) => {
				return renderSnippet(headers.runTime, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.runTime, row);
			}
		},
		{
			accessorKey: 'blockSize',
			header: ({ column }) => {
				return renderSnippet(headers.blockSize, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.blockSize, row);
			}
		},
		{
			accessorKey: 'fileSize',
			header: ({ column }) => {
				return renderSnippet(headers.fileSize, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.fileSize, row);
			}
		},
		{
			accessorKey: 'ioDepth',
			header: ({ column }) => {
				return renderSnippet(headers.ioDepth, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.ioDepth, row);
			}
		},
		{
			accessorKey: 'createdBy',
			header: ({ column }) => {
				return renderSnippet(headers.createdBy, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.createdBy, row);
			}
		},
		{
			accessorKey: 'bandwidth',
			header: ({ column }) => {
				return renderSnippet(headers.bandwidth, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.bandwidth, row);
			}
		},
		{
			accessorKey: 'iops',
			header: ({ column }) => {
				return renderSnippet(headers.iops, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.iops, row);
			}
		},
		{
			accessorKey: 'latencyMinimum',
			header: ({ column }) => {
				return renderSnippet(headers.latencyMinimum, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.latencyMinimum, row);
			}
		},
		{
			accessorKey: 'latencyMaximum',
			header: ({ column }) => {
				return renderSnippet(headers.latencyMaximum, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.latencyMaximum, row);
			}
		},
		{
			accessorKey: 'latencyMean',
			header: ({ column }) => {
				return renderSnippet(headers.latencyMean, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.latencyMean, row);
			}
		},
		{
			accessorKey: 'startedAt',
			header: ({ column }) => {
				return renderSnippet(headers.startedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.startedAt, row);
			}
		},
		{
			accessorKey: 'completedAt',
			header: ({ column }) => {
				return renderSnippet(headers.completedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.completedAt, row);
			}
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, { row, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
