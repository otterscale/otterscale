<script lang="ts" generics="TData">
	import type { Table } from '@tanstack/table-core';

	import { Badge } from '$lib/components/ui/badge';
	import { m } from '$lib/paraglide/messages';

	let { table }: { table: Table<TData> } = $props();

	const selected = $derived(table.getFilteredSelectedRowModel().rows.length);
	const filtered = $derived(table.getFilteredRowModel().rows.length);
	const rows = $derived(table.getCoreRowModel().rows.length);
</script>

<span class="flex flex-wrap gap-1">
	<Badge variant="outline">
		{selected}
		{m.selected()}
	</Badge>
	<Badge variant="outline">
		{filtered}
		{m.filtered()}
	</Badge>
	<Badge variant="outline">
		{rows}
		{m.rows()}
	</Badge>
</span>
