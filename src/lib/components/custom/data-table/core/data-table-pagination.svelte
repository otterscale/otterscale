<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { Button } from '$lib/components/ui/button/index.js';

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
		<p>{table.getPageCount() > 0 ? table.getState().pagination.pageIndex + 1 : 0}</p>
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
