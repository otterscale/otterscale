import type { OSD } from "$gen/api/storage/v1/storage_pb";
import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import DataTableActions from './actions.svelte';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<OSD>[] = [
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
        accessorKey: "stateUp",
        header: ({ column }) => {
            return renderSnippet(headers.stateUp, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.stateUp, row);
        },
    },
    {
        accessorKey: "stateIn",
        header: ({ column }) => {
            return renderSnippet(headers.stateIn, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.stateIn, row);
        },
    },
    {
        accessorKey: "exists",
        header: ({ column }) => {
            return renderSnippet(headers.exists, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.exists, row);
        },
    },
    {
        accessorKey: "deviceClass",
        header: ({ column }) => {
            return renderSnippet(headers.deviceClass, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.deviceClass, row);
        },
    },
    {
        accessorKey: "placementGroupCount",
        header: ({ column }) => {
            return renderSnippet(headers.placementGroupCount, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.placementGroupCount, row);
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
    {
        accessorKey: "readBytes",
        header: ({ column }) => {
            return renderSnippet(headers.readBytes, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.readBytes, row);
        },
    },
    {
        accessorKey: "writeBytes",
        header: ({ column }) => {
            return renderSnippet(headers.writeBytes, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.writeBytes, row);
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            return renderComponent(DataTableActions, { osd: row.original });
        },
    },
];

export {
    columns
};
