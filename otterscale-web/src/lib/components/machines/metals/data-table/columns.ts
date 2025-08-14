import type { Machine } from '$lib/api/machine/v1/machine_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Machine>[] = [
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
        accessorKey: "fqdn_ip",
        header: ({ column }) => {
            return renderSnippet(headers.fqdn_ip, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.fqdn_ip, row);
        },
    },
    {
        id: "fqdn",
        filterFn: (row, columnId, filterValue: string | undefined) => {
            if (filterValue === undefined) {
                return true
            }

            return row.original.fqdn.includes(filterValue)
        }
    },
    {
        accessorKey: "powerState",
        header: ({ column }) => {
            return renderSnippet(headers.powerState, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.powerState, row);
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
        accessorKey: "cores_arch",
        header: ({ column }) => {
            return renderSnippet(headers.cores_arch, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.cores_arch, row);
        },
    },
    {
        accessorKey: "ram",
        header: ({ column }) => {
            return renderSnippet(headers.ram, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.ram, row);
        },
    },
    {
        accessorKey: "disk",
        header: ({ column }) => {
            return renderSnippet(headers.disk, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.disk, row);
        },
    },
    {
        accessorKey: "storage",
        header: ({ column }) => {
            return renderSnippet(headers.storage, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.storage, row);
        },
    },
];

export {
    columns
};
