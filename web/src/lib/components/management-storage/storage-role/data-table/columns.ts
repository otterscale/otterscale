import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";

import DataTableActions from "./actions.svelte";
import { type Role } from './types';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Role>[] = [
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
        accessorKey: "roleName",
        header: ({ column }) => {
            return renderSnippet(headers.roleName, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.roleName, row);
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
        accessorKey: "arn",
        header: ({ column }) => {
            return renderSnippet(headers.arn, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.arn, row);
        },
    },
    {
        accessorKey: "createDate",
        header: ({ column }) => {
            return renderSnippet(headers.createDate, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createDate, row);
        },
    },
    {
        accessorKey: "maximumSessionDuration",
        header: ({ column }) => {
            return renderSnippet(headers.maximumSessionDuration, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.maximumSessionDuration, row);
        },
        filterFn: 'inNumberRange',
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { role: row.original });
        },
    },
];

export {
    columns
};
