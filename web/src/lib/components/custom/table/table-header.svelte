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

<Table.Header
	bind:ref
	data-slot="table-header"
	class={cn(
		'rounded-lg text-xs font-medium capitalize [&_th]:bg-muted [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</Table.Header>
