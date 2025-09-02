<script lang="ts" module>
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Table as TableType } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table, class: className }: { table: TableType<TData>; class?: string } = $props();
</script>

<Table.Row>
	<Table.Cell colspan={Object.keys(table.getHeaderGroups().flatMap((headerGroup) => headerGroup.headers)).length}>
		<div class={cn('flex items-center justify-center gap-1 p-4', className)}>
			<Icon icon="ph:magnifying-glass" />
			{m.no_data()}
		</div>
	</Table.Cell>
</Table.Row>
