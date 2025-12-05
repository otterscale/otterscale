<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { cn } from '$lib/utils';

	let { row }: { row: Row<TData> } = $props();
</script>

<div data-slot="data-table-row-expander" class="flex items-center justify-center">
	<button aria-label="Toggle Row Expansion" onclick={() => row.toggleExpanded()}>
		<Icon
			icon="ph:caret-down"
			class={cn(
				'size-3 transition-transform duration-200',
				row.getIsExpanded() ? '-rotate-180' : ''
			)}
		/>
	</button>
</div>
