<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Table } from '@tanstack/table-core';

	import { cn } from '$lib/utils';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();
</script>

<div data-slot="data-table-row-expander" class="flex items-center justify-center">
	<button aria-label="Toggle Row Expansion" onclick={() => table.toggleAllRowsExpanded()}>
		<Icon
			icon="ph:caret-down-fill"
			class={cn(
				'size-4 transition-transform duration-200',
				table.getIsAllRowsExpanded() ? '-rotate-180' : ''
			)}
		/>
	</button>
</div>
