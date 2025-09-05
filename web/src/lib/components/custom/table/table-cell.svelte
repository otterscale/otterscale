<script lang="ts" module>
	import type { HTMLTdAttributes } from 'svelte/elements';

	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
</script>

<script lang="ts">
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
