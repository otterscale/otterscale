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

<Table.Body
	bind:ref
	data-slot="table-body"
	class={cn('[&_tr]:last:rounded-b-lg', className)}
	{...restProps}
>
	{@render children?.()}
</Table.Body>
