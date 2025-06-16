<script lang="ts" generics="TData">
	import { Badge } from '$lib/components/ui/badge';
	import type { Table } from '@tanstack/table-core';

	let { table, data }: { table: Table<TData>; data: TData[] } = $props();

	const selected = $derived(table.getFilteredSelectedRowModel().rows.length);
	const filtered = $derived(table.getFilteredRowModel().rows.length);
	const rows = $derived(table.getCoreRowModel().rows.length);
</script>

<span class=" flex flex-wrap gap-1">
	<Badge variant="outline">
		{selected} selected
	</Badge>
	<Badge variant="outline">
		{filtered} filtered
	</Badge>
	<Badge variant="outline">
		{rows} rows
	</Badge>
</span>
