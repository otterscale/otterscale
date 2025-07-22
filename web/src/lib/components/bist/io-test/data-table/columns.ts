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
        accessorKey: "accessMode",
        header: ({ column }) => {
            return renderSnippet(headers.accessMode, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.accessMode, row);
        },
    },
    {
        accessorKey: "jobCount",
        header: ({ column }) => {
            return renderSnippet(headers.jobCount, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.jobCount, row);
        },
    },
    {
        accessorKey: "runTime",
        header: ({ column }) => {
            return renderSnippet(headers.runTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.runTime, row);
        },
    },
    {
        accessorKey: "blockSize",
        header: ({ column }) => {
            return renderSnippet(headers.blockSize, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.blockSize, row);
        },
    },
    {
        accessorKey: "fileSize",
        header: ({ column }) => {
            return renderSnippet(headers.fileSize, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.fileSize, row);
        },
    },
    {
        accessorKey: "ioDepth",
        header: ({ column }) => {
            return renderSnippet(headers.ioDepth, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.ioDepth, row);
        },
    },
    {
        accessorKey: "groupID",
        header: ({ column }) => {
            return renderSnippet(headers.groupID, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.groupID, row);
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
