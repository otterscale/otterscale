import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import type { Bucket } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<Bucket>[] = [
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
        accessorKey: "owner",
        header: ({ column }) => {
            return renderSnippet(headers.owner, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.owner, row);
        },
    },
    {
        accessorKey: "usedCapacity",
        header: ({ column }) => {
            return renderSnippet(headers.usedCapacity, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.usedCapacity, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        accessorKey: "capacityLimit",
        header: ({ column }) => {
            return renderSnippet(headers.capacityLimit, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.capacityLimit, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        accessorKey: "objects",
        header: ({ column }) => {
            return renderSnippet(headers.objects, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.objects, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        accessorKey: "objectLimit",
        header: ({ column }) => {
            return renderSnippet(headers.objectLimit, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.objectLimit, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { bucket: row.original });
        },
    },
];

export {
    columns
}
