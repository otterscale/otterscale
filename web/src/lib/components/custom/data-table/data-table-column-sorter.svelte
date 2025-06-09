<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { Column } from '@tanstack/table-core';

	let { column }: { column: Column<TData> } = $props();
</script>

<Button variant="ghost" onclick={column.getToggleSortingHandler()}>
	{#if column.getIsSorted() == 'asc'}
		<Icon icon="ph:sort-ascending" />
	{:else if column.getIsSorted() == 'desc'}
		<Icon icon="ph:sort-descending" />
	{:else}
		<Icon icon="ph:caret-up-down" />
	{/if}
</Button>
