<script lang="ts">
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import * as Table from '$lib/components/ui/table';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLTableRowElement>> = $props();
</script>

<Table.Row bind:ref data-slot="table-row" class={cn(className)} {...restProps}>
	{@render children?.()}
</Table.Row>
