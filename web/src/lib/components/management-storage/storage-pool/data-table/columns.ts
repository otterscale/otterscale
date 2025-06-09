import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import { type Pool } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<Pool>[] = [
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
        accessorKey: "dataProtection",
        header: ({ column }) => {
            return renderSnippet(headers.dataProtection, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.dataProtection, row);
        },
    },
    {
        accessorKey: "applications",
        header: ({ column }) => {
            return renderSnippet(headers.applications, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.applications, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "PGStatus",
        header: ({ column }) => {
            return renderSnippet(headers.PGStatus, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.PGStatus, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        accessorKey: "usage",
        header: ({ column }) => {
            return renderSnippet(headers.usage, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.usage, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { pool: row.original });
        },
    },
];

export {
    columns
}