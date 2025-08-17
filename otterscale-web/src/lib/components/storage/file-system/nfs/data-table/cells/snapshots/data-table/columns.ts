import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Subvolume_Snapshot>[] = [
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
        accessorKey: "hasPendingClones",
        header: ({ column }) => {
            return renderSnippet(headers.hasPendingClones, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.hasPendingClones, row);
        },
    },
    {
        accessorKey: "actions",
        header: ({ column }) => {
            return renderSnippet(headers.actions, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.actions, row);
        },
    },
];

export {
    columns
};
