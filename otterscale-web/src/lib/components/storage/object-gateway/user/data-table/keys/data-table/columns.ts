import type { User_Key } from '$lib/api/storage/v1/storage_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<User_Key>[] = [
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
        accessorKey: "access",
        header: ({ column }) => {
            return renderSnippet(headers.access, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.access, row);
        },
    },
];

export {
    columns
};
