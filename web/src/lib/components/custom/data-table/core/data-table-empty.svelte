<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import type { Table as TableType } from '@tanstack/table-core';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { cn } from '$lib/utils';

	let { table, class: className }: { table: TableType<TData>; class?: string } = $props();
</script>

<Table.Row>
	<Table.Cell
		colspan={Object.keys(table.getHeaderGroups().flatMap((headerGroup) => headerGroup.headers))
			.length}
	>
		<div class={cn('flex items-center justify-center gap-1 p-8', className)}>
			<Empty.Root>
				<Empty.Header>
					<Empty.Media variant="icon">
						<Icon icon="ph:magnifying-glass" />
					</Empty.Media>
					<Empty.Title>No data</Empty.Title>
					<Empty.Description>
						Please adjust your filters or add new entries to see data here.
					</Empty.Description>
				</Empty.Header>
			</Empty.Root>
		</div>
	</Table.Cell>
</Table.Row>
