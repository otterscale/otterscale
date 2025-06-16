import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";
import type { Namespace } from './types'
import DataTableActions from "./actions.svelte";
import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<Namespace>[] = [
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
        accessorKey: "pool",
        header: ({ column }) => {
            return renderSnippet(headers.pool, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.pool, row);
        },
    },
    {
        accessorKey: "namespace",
        header: ({ column }) => {
            return renderSnippet(headers.namespace, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.namespace, row);
        },
    },
    {
        accessorKey: "totalImages",
        header: ({ column }) => {
            return renderSnippet(headers.totalImages, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.totalImages, row);
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { namespace: row.original });
        },
    },
];

export {
    columns
}