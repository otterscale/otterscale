import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { type Subvolume } from './types';

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
        accessorKey: "dataPool",
        header: ({ column }) => {
            return renderSnippet(headers.dataPool, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.dataPool, row);
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
        accessorKey: "path",
        header: ({ column }) => {
            return renderSnippet(headers.path, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.path, row);
        },
    },
    {
        accessorKey: "mode",
        header: ({ column }) => {
            return renderSnippet(headers.mode, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.mode, row);
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
            return renderComponent(DataTableActions, { subvolume: row.original });
        },
    },
];

export {
    columns
};
