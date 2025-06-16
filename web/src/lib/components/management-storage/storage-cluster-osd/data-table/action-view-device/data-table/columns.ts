import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { Device } from './types'
import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<Device>[] = [
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
        accessorKey: "name",
        header: ({ column }) => {
            return renderSnippet(headers.name, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.name, row);
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
        filterFn: "arrIncludesSome"
    },
    {
        accessorKey: "lifeExpectancy",
        header: ({ column }) => {
            return renderSnippet(headers.lifeExpectancy, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.lifeExpectancy, row);
        },
    },
    {
        accessorKey: "daemons",
        header: ({ column }) => {
            return renderSnippet(headers.daemons, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.daemons, row);
        },
    },
]

export {
    columns
}