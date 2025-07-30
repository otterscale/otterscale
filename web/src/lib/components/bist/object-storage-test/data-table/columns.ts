import { type TestResult } from '$gen/api/bist/v1/bist_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import { type ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<TestResult>[] = [
    {
        accessorKey: "name",
        header: ({ column }) => {
            return renderSnippet(headers.name, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.name, row);
        },
    },
    {
        accessorKey: "status",
        header: ({ column }) => {
            return renderSnippet(headers.status, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.status, row);
        },
    },
    {
        accessorKey: "target",
        header: ({ column }) => {
            return renderSnippet(headers.target, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.target, row);
        },
    },
    {
        accessorKey: "createdBy",
        header: ({ column }) => {
            return renderSnippet(headers.createdBy, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createdBy, row);
        },
    },
    {
        accessorKey: "operation",
        header: ({ column }) => {
            return renderSnippet(headers.operation, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.operation, row);
        },
    },
    {
        accessorKey: "duration",
        header: ({ column }) => {
            return renderSnippet(headers.duration, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.duration, row);
        },
    },
    {
        accessorKey: "objectSize",
        header: ({ column }) => {
            return renderSnippet(headers.objectSize, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.objectSize, row);
        },
    },
    {
        accessorKey: "objectCount",
        header: ({ column }) => {
            return renderSnippet(headers.objectCount, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.objectCount, row);
        },
    },
    {
        accessorKey: "throughputFastest",
        header: ({ column }) => {
            return renderSnippet(headers.throughputFastest, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.throughputFastest, row);
        },
    },
    {
        accessorKey: "throughputSlowest",
        header: ({ column }) => {
            return renderSnippet(headers.throughputSlowest, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.throughputSlowest, row);
        },
    },
    {
        accessorKey: "throughputMedian",
        header: ({ column }) => {
            return renderSnippet(headers.throughputMedian, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.throughputMedian, row);
        },
    },
    {
        accessorKey: "startedAt",
        header: ({ column }) => {
            return renderSnippet(headers.startedAt, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.startedAt, row);
        },
    },
    {
        accessorKey: "completedAt",
        header: ({ column }) => {
            return renderSnippet(headers.completedAt, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.completedAt, row);
        },
    }
];

export {
    columns
};
