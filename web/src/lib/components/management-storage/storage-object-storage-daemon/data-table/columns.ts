import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { type ObjectStorageDaemon } from './types';

const columns: ColumnDef<ObjectStorageDaemon>[] = [
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
        accessorKey: "id",
        header: ({ column }) => {
            return renderSnippet(headers.id, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.id, row);
        },
    },
    {
        accessorKey: "host",
        header: ({ column }) => {
            return renderSnippet(headers.host, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.host, row);
        },
    },
    {
        accessorKey: "devices",
        header: ({ column }) => {
            return renderSnippet(headers.devices, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.devices, row);
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
        filterFn: "arrIncludesSome"
    },
    {
        accessorKey: "deviceClass",
        header: ({ column }) => {
            return renderSnippet(headers.deviceClass, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.deviceClass, row);
        },
    },
    {
        accessorKey: "pgs",
        header: ({ column }) => {
            return renderSnippet(headers.pgs, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.pgs, row);
        },
    },
    {
        accessorKey: "size",
        header: ({ column }) => {
            return renderSnippet(headers.size, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.size, row);
        },
    },
    {
        accessorKey: "flags",
        header: ({ column }) => {
            return renderSnippet(headers.flags, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.flags, row);
        },
        filterFn: "arrIncludesSome"
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
        accessorKey: "readBytes",
        header: ({ column }) => {
            return renderSnippet(headers.readBytes, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.readBytes, row);
        },
    },
    {
        accessorKey: "writeBytes",
        header: ({ column }) => {
            return renderSnippet(headers.writeBytes, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.writeBytes, row);
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { osd: row.original });
        },
    },
];

export {
    columns
};
