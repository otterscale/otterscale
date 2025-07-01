import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { MON } from "$gen/api/storage/v1/storage_pb";

const columns: ColumnDef<MON>[] = [
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
        accessorKey: "leader",
        header: ({ column }) => {
            return renderSnippet(headers.leader, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.leader, row);
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
        accessorKey: "rank",
        header: ({ column }) => {
            return renderSnippet(headers.rank, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.rank, row);
        },
    },
    {
        accessorKey: "publicAddress",
        header: ({ column }) => {
            return renderSnippet(headers.publicAddress, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.publicAddress, row);
        },
    },
];

export {
    columns
};
