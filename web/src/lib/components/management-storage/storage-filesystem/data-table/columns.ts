import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import { type FileSystem } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<FileSystem>[] = [
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
        accessorKey: "enabled",
        header: ({ column }) => {
            return renderSnippet(headers.enabled, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.enabled, row);
        },
    },
    {
        accessorKey: "permission",
        header: ({ column }) => {
            return renderSnippet(headers.permission, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.permission, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "createTime",
        header: ({ column }) => {
            return renderSnippet(headers.createTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createTime, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { fileSystem: row.original });
        },
    },
];

export {
    columns
}