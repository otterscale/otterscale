<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLTdAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLTdAttributes> = $props();
</script>

<Table.Cell bind:ref data-slot="table-cell" class={cn('text-xs', className)} {...restProps}>
	{@render children?.()}
</Table.Cell>
