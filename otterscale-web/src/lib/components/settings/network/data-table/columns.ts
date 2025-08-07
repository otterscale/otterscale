import type { Network } from '$lib/api/network/v1/network_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const columns: ColumnDef<Network>[] = [
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
        accessorKey: "fabric",
        header: ({ column }) => {
            return renderSnippet(headers.fabric, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.fabric, row);
        },
    },
    {
        id: "fabricName",
        filterFn: (row, columnId, filterValue) => {
            if (!filterValue) {
                return true
            }

            return row.original.fabric?.name === filterValue
        },
    },
    {
        accessorKey: "vlan",
        header: ({ column }) => {
            return renderSnippet(headers.vlan, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.vlan, row);
        },
    },
    {
        id: "vlanName",
        filterFn: (row, columnId, filterValue) => {
            console.log(row.original.vlan?.name, filterValue)
            if (!filterValue) {
                return true
            }

            return filterValue.includes(row.original.vlan?.name)
        },
    },
    {
        accessorKey: "dhcpOn",
        header: ({ column }) => {
            return renderSnippet(headers.dhcpOn, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.dhcpOn, row);
        },
    },
    {
        accessorKey: "subnet",
        header: ({ column }) => {
            return renderSnippet(headers.subnet, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.subnet, row);
        },
    },
    {
        accessorKey: "ipAddresses",
        header: ({ column }) => {
            return renderSnippet(headers.ipAddresses, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.ipAddresses, row);
        },
    },
    {
        accessorKey: "ipRanges",
        header: ({ column }) => {
            return renderSnippet(headers.ipRanges, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.ipRanges, row);
        },
    },
    {
        accessorKey: "statistics",
        header: ({ column }) => {
            return renderSnippet(headers.statistics, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.statistics, row);
        },
    },
];

export {
    columns
};
