import { type TestResult } from '$gen/api/bist/v1/bist_pb';
import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import { type ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<TestResult>[] = [
    {
        accessorKey: "uid",
        header: ({ column }) => {
            return renderSnippet(headers.uid, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.uid, row);
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
        accessorKey: "status",
        header: ({ column }) => {
            return renderSnippet(headers.status, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.status, row);
        },
    },
    {
        accessorKey: "createdBy",
        header: ({ column }) => {
            return renderSnippet(headers.createdBy, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createdBy, row);
        },
    },
    {
        accessorKey: "startedAt",
        header: ({ column }) => {
            return renderSnippet(headers.startedAt, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.startedAt, row);
        },
    },
    {
        accessorKey: "completedAt",
        header: ({ column }) => {
            return renderSnippet(headers.completedAt, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.completedAt, row);
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { flexibleIOTest: row.original });
        },
    },
];

export {
    columns
};
