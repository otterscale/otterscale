<script lang="ts" module>
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import type { Column } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { column }: { column: Column<TData> } = $props();
</script>

<Button
	data-slot="data-table-sorter"
	variant="ghost"
	onclick={column.getToggleSortingHandler()}
	class="invisible group-hover:visible"
>
	{#if column.getIsSorted() == 'asc'}
		<Icon icon="ph:sort-ascending" />
	{:else if column.getIsSorted() == 'desc'}
		<Icon icon="ph:sort-descending" />
	{:else}
		<Icon icon="ph:list" />
	{/if}
</Button>
