import type { Bucket } from "$gen/api/storage/v1/storage_pb";
import { renderSnippet } from "$lib/components/ui/data-table/index.js";
import type { ColumnDef } from "@tanstack/table-core";
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { timestampDate, type Timestamp } from "@bufbuild/protobuf/wkt";

const columns: ColumnDef<Bucket>[] = [
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
        accessorKey: "owner",
        header: ({ column }) => {
            return renderSnippet(headers.owner, column)
        },
        cell: ({ row }) => {
            return renderSnippet(cells.owner, row);
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

