<script lang="ts" module>
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLTableSectionElement>> = $props();
</script>

<Table.Footer bind:ref data-slot="table-footer" class={cn(className)} {...restProps}>
	{@render children?.()}
</Table.Footer>
