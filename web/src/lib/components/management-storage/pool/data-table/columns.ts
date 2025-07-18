import type { Pool } from '$gen/api/storage/v1/storage_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Pool>[] = [
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
        accessorKey: "poolType",
        header: ({ column }) => {
            return renderSnippet(headers.type, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.type, row);
        },
    },
    {
        accessorKey: "applications",
        header: ({ column }) => {
            return renderSnippet(headers.applications, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.applications, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "placementGroupState",
        header: ({ column }) => {
            return renderSnippet(headers.placement_group_state, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.placement_group_state, row);
        },
        filterFn: (row, columnId, filterValue) => {
            const value = Object.keys(row.getValue(columnId) ?? {});
            if (!value.length || !filterValue.length) return true;
            return value.some(v => filterValue.includes(v));
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
    // {
    //     accessorKey: "readBytes",
    //     header: ({ column }) => {
    //         return renderSnippet(headers.readBytes, column)
    //     },
    //     cell: ({ row }) => {
    //         return renderSnippet(cells.readBytes, row);
    //     },
    // },
    // {
    //     accessorKey: "writeBytes",
    //     header: ({ column }) => {
    //         return renderSnippet(headers.writeBytes, column)
    //     },
    //     cell: ({ row }) => {
    //         return renderSnippet(cells.writeBytes, row);
    //     },
    // },
];

export {
    columns
};
