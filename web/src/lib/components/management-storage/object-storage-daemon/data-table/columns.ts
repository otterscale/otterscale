import type { OSD } from "$gen/api/storage/v1/storage_pb";
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<OSD>[] = [
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
        accessorKey: "up",
        header: ({ column }) => {
            return renderSnippet(headers.stateUp, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.stateUp, row);
        },
    },
    {
        accessorKey: "in",
        header: ({ column }) => {
            return renderSnippet(headers.stateIn, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.stateIn, row);
        },
    },
    {
        accessorKey: "exists",
        header: ({ column }) => {
            return renderSnippet(headers.exists, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.exists, row);
        },
        filterFn: 'equals',
    },
    {
        accessorKey: "deviceClass",
        header: ({ column }) => {
            return renderSnippet(headers.deviceClass, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.deviceClass, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "machine",
        header: ({ column }) => {
            return renderSnippet(headers.machine, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.machine, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "placementGroupCount",
        header: ({ column }) => {
            return renderSnippet(headers.placementGroupCount, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.placementGroupCount, row);
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
];

export {
    columns
};
