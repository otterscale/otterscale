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

<Table.Body bind:ref data-slot="table-body" class={cn(className)} {...restProps}>
	{@render children?.()}
</Table.Body>
