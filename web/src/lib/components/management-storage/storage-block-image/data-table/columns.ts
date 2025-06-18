import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { BlockImage } from './types';

const columns: ColumnDef<BlockImage>[] = [
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
        accessorKey: "usage",
        header: ({ column }) => {
            return renderSnippet(headers.usage, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.usage, row);
        },
    },

    {
        accessorKey: "objects",
        header: ({ column }) => {
            return renderSnippet(headers.objects, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.objects, row);
        },
    },
    {
        accessorKey: "parent",
        header: ({ column }) => {
            return renderSnippet(headers.parent, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.parent, row);
        },
    },
    {
        accessorKey: "mirroring",
        header: ({ column }) => {
            return renderSnippet(headers.mirroring, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.mirroring, row);
        },
    },
    {
        accessorKey: "nextScheduledSnapshot",
        header: ({ column }) => {
            return renderSnippet(headers.nextScheduledSnapshot, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.nextScheduledSnapshot, row);
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { blockImage: row.original });
        },
    },
];

export {
    columns
};
