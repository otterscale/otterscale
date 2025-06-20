import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";
import { type Snapshot } from './types'
import DataTableActions from "./actions.svelte";
import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<Snapshot>[] = [
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
        accessorKey: "createTime",
        header: ({ column }) => {
            return renderSnippet(headers.createTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createTime, row);
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
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { snapshot: row.original });
        },
    },
];

export {
    columns
}