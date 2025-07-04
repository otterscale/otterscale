import { type TestResult } from '$gen/api/bist/v1/bist_pb';
import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import { type ColumnDef } from "@tanstack/table-core";
import DataTableActions from "./actions.svelte";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<TestResult>[] = [
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
        accessorKey: "name",
        header: ({ column }) => {
            return renderSnippet(headers.name, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.name, row);
        },
    },
    {
        accessorKey: "input",
        header: ({ column }) => {
            return renderSnippet(headers.input, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.input, row);
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
