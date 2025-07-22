import type { Subvolume } from '$gen/api/storage/v1/storage_pb';
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { Row } from '$lib/components/ui/table';
import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';

const columns: ColumnDef<Subvolume>[] = [
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
        accessorKey: "poolName",
        header: ({ column }) => {
            return renderSnippet(headers.poolName, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.poolName, row);
        },
    },
    // {
    //     accessorKey: "mode",
    //     header: ({ column }) => {
    //         return renderSnippet(headers.mode, column)
    //     },
    //     cell: ({ row }) => {
    //         return renderSnippet(cells.mode, row);
    //     },
    // },
    {
        accessorKey: "export",
        header: ({ column }) => {
            return renderSnippet(headers.Export, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.Export, row);
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
        accessorKey: "createTime",
        header: ({ column }) => {
            return renderSnippet(headers.createTime, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.createTime, row);
        },
        sortingFn: (previousRow, nextRow, columnId) => {
            const previous: Timestamp | undefined = previousRow.original.createdAt
            const next: Timestamp | undefined = nextRow.original.createdAt

            if (!(previous || next)) {
                return 0
            }
            else if (!previous) {
                return -1
            }
            else if (!next) {
                return 1
            }
            else {
                if (timestampDate(previous) < timestampDate(next)) {
                    return -1
                }
                else if (timestampDate(previous) > timestampDate(next)) {
                    return 1
                }
                else {
                    return 0
                }
            }
        }
    },
];

export {
    columns
};
