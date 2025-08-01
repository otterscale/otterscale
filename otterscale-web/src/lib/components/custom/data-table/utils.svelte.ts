import {
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    type ColumnDef,
    type ColumnFiltersState,
    type PaginationState,
    type RowSelectionState,
    type SortingState,
    type Updater,
    type VisibilityState
} from '@tanstack/table-core';
import type { TableState } from './type';
// import type { PrometheusDriver } from 'prometheus-query';

// class StatisticManager<TData> {
//     table: Table<TData>

//     constructor(table: Table<TData>) {
//         this.table = table
//     }

//     get filteredData() {
//         return this.table.getFilteredRowModel().rows.map((row) => row.original)
//     }

//     count(key: string) {
//         const values = this.filteredData.map((datum) => datum[key as keyof TData])
//         const distinctValues: any[] = Array(...new Set(values));
//         return distinctValues.length;
//     }

//     groupCount(key: string) {
//         const values: number[] = this.filteredData.map((datum) => datum[key as keyof TData] as number)
//         const counts: Record<any, number> = values.reduce(
//             (a, value) => {
//                 a[value] = (a[value] || 0) + 1;
//                 return a;
//             },
//             {} as Record<any, number>
//         );
//         return Object.fromEntries(Object.entries(counts).sort(([, p], [, n]) => n - p));
//     }

//     sum(key: string) {
//         const values: number[] = this.filteredData.map((datum) => datum[key as keyof TData] as number)
//         return values.reduce((a, value) => a + value, 0);
//     }
// }

function getTableOptions(data: any[], columns: ColumnDef<any>[], tableState: TableState) {
    return {
        get data() {
            return data;
        },

        columns,

        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),

        state: {
            get pagination() {
                return tableState.pagination;
            },
            get sorting() {
                return tableState.sorting;
            },
            get columnFilters() {
                return tableState.columnFilters;
            },
            get columnVisibility() {
                return tableState.columnVisibility;
            },
            get rowSelection() {
                return tableState.rowSelection;
            }
        },

        onPaginationChange: (updater: Updater<PaginationState>) => {
            if (typeof updater === 'function') {
                tableState.pagination = updater(tableState.pagination);
            } else {
                tableState.pagination = updater;
            }
        },
        onSortingChange: (updater: Updater<SortingState>) => {
            if (typeof updater === 'function') {
                tableState.sorting = updater(tableState.sorting);
            } else {
                tableState.sorting = updater;
            }
        },
        onColumnFiltersChange: (updater: Updater<ColumnFiltersState>) => {
            if (typeof updater === 'function') {
                tableState.columnFilters = updater(tableState.columnFilters);
            } else {
                tableState.columnFilters = updater;
            }
        },
        onColumnVisibilityChange: (updater: Updater<VisibilityState>) => {
            if (typeof updater === 'function') {
                tableState.columnVisibility = updater(tableState.columnVisibility);
            } else {
                tableState.columnVisibility = updater;
            }
        },
        onRowSelectionChange: (updater: Updater<RowSelectionState>) => {
            if (typeof updater === 'function') {
                tableState.rowSelection = updater(tableState.rowSelection);
            } else {
                tableState.rowSelection = updater;
            }
        },

        autoResetAll: false
    };
}

export { getTableOptions };

