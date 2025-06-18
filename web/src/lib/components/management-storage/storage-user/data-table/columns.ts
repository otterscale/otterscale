import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { type User } from './types';

const columns: ColumnDef<User>[] = [
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
        accessorKey: "username",
        header: ({ column }) => {
            return renderSnippet(headers.username, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.username, row);
        },
    },
    {
        accessorKey: "tenant",
        header: ({ column }) => {
            return renderSnippet(headers.tenant, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.tenant, row);
        },
    },
    {
        accessorKey: "fullName",
        header: ({ column }) => {
            return renderSnippet(headers.fullName, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.fullName, row);
        },
    },
    {
        accessorKey: "emailAddress",
        header: ({ column }) => {
            return renderSnippet(headers.emailAddress, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.emailAddress, row);
        },
    },
    {
        accessorKey: "suspended",
        header: ({ column }) => {
            return renderSnippet(headers.suspended, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.suspended, row);
        },
    },
    {
        accessorKey: "maximumBuckets",
        header: ({ column }) => {
            return renderSnippet(headers.maximumBuckets, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.maximumBuckets, row);
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
            return renderComponent(DataTableActions, { user: row.original });
        },
    },
];

export {
    columns
};
