<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import type { Table } from '@tanstack/table-core';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const selected = $derived(table.getFilteredSelectedRowModel().rows.length);
	const filtered = $derived(table.getFilteredRowModel().rows.length);
	const rows = $derived(table.getCoreRowModel().rows.length);
</script>

<span class="flex flex-wrap gap-1">
	<Badge variant="outline">
		{selected}
		{m.datatable_footer_selected()}
	</Badge>
	<Badge variant="outline">
		{filtered}
		{m.datatable_footer_filtered()}
	</Badge>
	<Badge variant="outline">
		{rows}
		{m.datatable_footer_rows()}
	</Badge>
</span>
