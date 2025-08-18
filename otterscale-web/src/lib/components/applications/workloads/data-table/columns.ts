import type { Application } from '$lib/api/application/v1/application_pb';
import { getSortingFunction } from '$lib/components/custom/data-table';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Application>[] = [
    {
        id: "select",
        header: ({ table }) => {
            return renderSnippet(headers.row_picker, table)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.row_picker, row);
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
        filterFn: 'arrIncludesSome'
    },
    {
        accessorKey: "namespace",
        header: ({ column }) => {
            return renderSnippet(headers.namespace, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.namespace, row);
        },
        filterFn: 'arrIncludesSome'
    },
    {
        accessorKey: "health",
        header: ({ column }) => {
            return renderSnippet(headers.health, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.health, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.healthies / previousRow.original.pods.length,
                nextRow.original.healthies / nextRow.original.pods.length,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
    },
    {
        accessorKey: "service",
        header: ({ column }) => {
            return renderSnippet(headers.service, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.service, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.services.length,
                nextRow.original.services.length,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
    },
    {
        accessorKey: "pod",
        header: ({ column }) => {
            return renderSnippet(headers.pod, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.pod, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.pods.length,
                nextRow.original.pods.length,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
    },
    {
        accessorKey: "replica",
        header: ({ column }) => {
            return renderSnippet(headers.replica, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.replica, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.replicas,
                nextRow.original.replicas,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
    },
    {
        accessorKey: "container",
        header: ({ column }) => {
            return renderSnippet(headers.container, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.container, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.containers.length,
                nextRow.original.containers.length,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
    },
    {
        accessorKey: "volume",
        header: ({ column }) => {
            return renderSnippet(headers.volume, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.volume, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => (
            getSortingFunction(
                previousRow.original.persistentVolumeClaims.length,
                nextRow.original.persistentVolumeClaims.length,
                (p, n) => (p < n),
                (p, n) => (p === n)
            )
        )
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
