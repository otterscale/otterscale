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
	}: WithElementRef<HTMLAttributes<HTMLTableRowElement>> = $props();
</script>

<Table.Row
	bind:ref
	data-slot="table-row"
	class={cn(
		'hover:[&,&>svelte-css-wrapper]:[&>th,td]:bg-transparent',
		'text-xs hover:bg-muted/30 last-of-type:[&_td]:first:rounded-bl-lg last-of-type:[&_td]:last:rounded-br-lg',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</Table.Row>
