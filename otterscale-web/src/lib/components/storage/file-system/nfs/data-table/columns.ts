import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/utils.svelte';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Subvolume>[] = [
    {
        id: "select",
        header: ({ table }) => {
            return renderSnippet(headers._row_picker, table)
        },
        cell: ({ row }) => {
            return renderSnippet(cells._row_picker, row);
        },
        enableSorting: false,
        enableHiding: false,
    },
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
        accessorKey: "poolName",
        header: ({ column }) => {
            return renderSnippet(headers.poolName, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.poolName, row);
        },
    },
    {
        accessorKey: "export",
        header: ({ column }) => {
            return renderSnippet(headers.Export, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.Export, row);
        },
    },
    {
        accessorKey: "usage",
        header: ({ column }) => {
            return renderSnippet(headers.usage, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.usage, row);
        },
    },
    {
        accessorKey: "createTime",
        header: ({ column }) => {
            return renderSnippet(headers.createTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createTime, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.createdAt,
                nextRow.original.createdAt,
                (p, n) => (timestampDate(p) < timestampDate(n)),
                (p, n) => (timestampDate(p) === timestampDate(n))
            )
        )
    },
    {
        accessorKey: "snapshots",
        header: ({ column }) => {
            return renderSnippet(headers.snapshots, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.snapshots, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.snapshots.length,
                nextRow.original.snapshots.length,
                (p, n) => (timestampDate(p) < timestampDate(n)),
                (p, n) => (timestampDate(p) === timestampDate(n))
            )
        )
    },
    {
        accessorKey: "actions",
        header: ({ column }) => {
            return renderSnippet(headers.actions, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.actions, row);
        },
        
    },
];

export {
    columns
};
