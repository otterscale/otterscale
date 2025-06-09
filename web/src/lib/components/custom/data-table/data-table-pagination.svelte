<script lang="ts" generics="TData">
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';

	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TData> } = $props();
</script>

<div class="flex items-center justify-end gap-2">
	<Button
		variant="outline"
		size="icon"
		onclick={() => table.firstPage()}
		disabled={!table.getCanPreviousPage()}
	>
		<Icon icon="ph:caret-line-left" />
	</Button>
	<Button
		variant="outline"
		size="icon"
		onclick={() => table.previousPage()}
		disabled={!table.getCanPreviousPage()}
	>
		<Icon icon="ph:caret-left" />
	</Button>
	<span class="flex items-center gap-1 text-xs font-light">
		<p>{table.getState().pagination.pageIndex + 1}</p>
		<p>/</p>
		<p>{table.getPageCount()}</p>
	</span>
	<Button
		variant="outline"
		size="icon"
		onclick={() => table.nextPage()}
		disabled={!table.getCanNextPage()}
	>
		<Icon icon="ph:caret-right" />
	</Button>
	<Button
		variant="outline"
		size="icon"
		onclick={() => table.lastPage()}
		disabled={!table.getCanNextPage()}
	>
		<Icon icon="ph:caret-line-right" />
	</Button>
</div>
