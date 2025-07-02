import type { ColumnDef } from "@tanstack/table-core";
import { renderSnippet, renderComponent } from "$lib/components/ui/data-table/index.js";

import { type FlexibleIOTest } from './types'
import DataTableActions from "./actions.svelte";

import { headers } from './headers.svelte'
import { cells } from './cells.svelte'

const columns: ColumnDef<FlexibleIOTest>[] = [
    // {
    //     id: "select",
    //     header: ({ table }) => {
    //         return renderSnippet(headers._row_picker, table)
    //     },
    //     cell: ({ row }) => {
    //         return renderSnippet(cells._row_picker, row);
    //     },
    //     enableSorting: false,
    //     enableHiding: false,
    // },
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
        accessorKey: "rwMode",
        header: ({ column }) => {
            return renderSnippet(headers.rwMode, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.rwMode, row);
        },
    },
    {
        accessorKey: "fileSize",
        header: ({ column }) => {
            return renderSnippet(headers.fileSize, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.fileSize, row);
        },
        filterFn: 'arrIncludesSome',
    },
    {
        accessorKey: "numberJobs",
        header: ({ column }) => {
            return renderSnippet(headers.numberJobs, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.numberJobs, row);
        },
    },
    {
        accessorKey: "blockSize",
        header: ({ column }) => {
            return renderSnippet(headers.blockSize, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.blockSize, row);
        },
    },
    {
        accessorKey: "runtime",
        header: ({ column }) => {
            return renderSnippet(headers.runtime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.runtime, row);
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
        accessorKey: "modifyTime",
        header: ({ column }) => {
            return renderSnippet(headers.modifyTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.modifyTime, row);
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
}