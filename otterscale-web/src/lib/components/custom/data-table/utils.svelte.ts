import type { Table } from '@tanstack/table-core';

class StatisticManager<TData> {
    table: Table<TData>

    constructor(table: Table<TData>) {
        this.table = table
    }

    get filteredData() {
        return this.table.getFilteredRowModel().rows.map((row) => row.original)
    }

    count(key: string) {
        const values = this.filteredData.map((datum) => datum[key as keyof TData])
        const distinctValues: any[] = Array(...new Set(values));
        return distinctValues.length;
    }

    groupCount(key: string) {
        const values: number[] = this.filteredData.map((datum) => datum[key as keyof TData] as number)
        const counts: Record<any, number> = values.reduce(
            (a, value) => {
                a[value] = (a[value] || 0) + 1;
                return a;
            },
            {} as Record<any, number>
        );
        return Object.fromEntries(Object.entries(counts).sort(([, p], [, n]) => n - p));
    }

    sum(key: string) {
        const values: number[] = this.filteredData.map((datum) => datum[key as keyof TData] as number)
        return values.reduce((a, value) => a + value, 0);
    }
}

export {
    StatisticManager
}
