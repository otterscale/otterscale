<script lang="ts">
	import type { HTMLTdAttributes } from 'svelte/elements';

	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		alignClass,
		...restProps
	}: WithElementRef<HTMLTdAttributes> & {
		alignClass?: string;
	} = $props();
</script>

<Table.Cell bind:ref data-slot="table-cell" class={cn('p-2', className)} {...restProps}>
	<div class={cn('flex flex-col gap-1', alignClass)}>
		{@render children?.()}
	</div>
</Table.Cell>
