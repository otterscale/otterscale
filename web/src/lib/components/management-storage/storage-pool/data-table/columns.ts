import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { type Pool } from './types';

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
            return renderComponent(DataTableActions, { pool: row.original });
        },
    },
];

export {
    columns
};
