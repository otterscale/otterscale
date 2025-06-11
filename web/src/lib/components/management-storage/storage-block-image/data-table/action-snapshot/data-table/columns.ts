import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import type { BlockImageSnapshot } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<BlockImageSnapshot>[] = [
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
        accessorKey: "size",
        header: ({ column }) => {
            return renderSnippet(headers.size, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.size, row);
        },
    },
    {
        accessorKey: "used",
        header: ({ column }) => {
            return renderSnippet(headers.used, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.used, row);
        },
    },
    {
        accessorKey: "state",
        header: ({ column }) => {
            return renderSnippet(headers.state, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.state, row);
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
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { blockImagesnapshot: row.original });
        },
    },
];

export {
    columns
}