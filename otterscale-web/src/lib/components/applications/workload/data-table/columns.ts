import type { Application } from '$lib/api/application/v1/application_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Application>[] = [
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
        accessorKey: "type",
        header: ({ column }) => {
            return renderSnippet(headers.type, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.type, row);
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
        accessorKey: "health",
        header: ({ column }) => {
            return renderSnippet(headers.health, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.health, row);
        },
    },
    {
        accessorKey: "service",
        header: ({ column }) => {
            return renderSnippet(headers.service, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.service, row);
        },
    },
    {
        accessorKey: "pod",
        header: ({ column }) => {
            return renderSnippet(headers.pod, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.pod, row);
        },
    },
    {
        accessorKey: "replica",
        header: ({ column }) => {
            return renderSnippet(headers.replica, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.replica, row);
        },
    },
    {
        accessorKey: "container",
        header: ({ column }) => {
            return renderSnippet(headers.container, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.container, row);
        },
    },
    {
        accessorKey: "volume",
        header: ({ column }) => {
            return renderSnippet(headers.volume, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.volume, row);
        },
    },
    {
        accessorKey: "nodeport",
        header: ({ column }) => {
            return renderSnippet(headers.nodeport, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.nodeport, row);
        },
    },
];

export {
    columns
};
