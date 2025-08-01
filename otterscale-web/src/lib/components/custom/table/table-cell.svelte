<script lang="ts" module>
	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLTdAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLTdAttributes> = $props();
</script>

<Table.Cell
	bind:ref
	data-slot="table-cell"
	class={cn(
		'has-data-[slot=badge]:gap-1 has-data-[slot=badge]:flex has-data-[slot=badge]:flex-wrap has-data-[slot=badge]:items-center text-xs',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</Table.Cell>
