<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLTableSectionElement>> = $props();
</script>

<Table.Header bind:ref data-slot="table-header" class={cn(className)} {...restProps}>
	{@render children?.()}
</Table.Header>
