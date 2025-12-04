<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import type { Column } from '@tanstack/table-core';

	import { Button } from '$lib/components/ui/button/index.js';

	let { column }: { column: Column<TData> } = $props();
</script>

<Button
	data-slot="data-table-sorter"
	variant="ghost"
	size="icon-sm"
	onclick={column.getToggleSortingHandler()}
	class="invisible p-1 group-hover:visible"
>
	{#if column.getIsSorted() == 'asc'}
		<Icon icon="ph:caret-up-bold" class="size-3" />
	{:else if column.getIsSorted() == 'desc'}
		<Icon icon="ph:caret-down-bold" class="size-3" />
	{:else}
		<Icon icon="ph:minus-bold" class="size-3" />
	{/if}
</Button>
