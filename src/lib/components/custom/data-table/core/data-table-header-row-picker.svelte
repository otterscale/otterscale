<script lang="ts" generics="TData">
	import type { Table } from '@tanstack/table-core';

	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	let { table }: { table: Table<TData> } = $props();
</script>

<div data-slot="data-table-row-picker" class="flex items-center justify-center">
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
</div>
