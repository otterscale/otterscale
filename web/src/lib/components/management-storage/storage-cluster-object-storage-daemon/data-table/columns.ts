import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import { type OSD } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

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
        accessorKey: "status",
        header: ({ column }) => {
            return renderSnippet(headers.status, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.status, row);
        },
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
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { osd: row.original });
        },
    },
];

export {
    columns
}