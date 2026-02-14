<script lang="ts">
	import type { HTMLTdAttributes } from 'svelte/elements';

	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLTdAttributes> & {} = $props();
</script>

<Table.Cell bind:ref data-slot="table-cell" class={cn('p-2', className)} {...restProps}>
	{@render children?.()}
</Table.Cell>
