<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLThAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLThAttributes> = $props();
</script>

<Table.Head
	bind:ref
	data-slot="table-head"
	class={cn('text-xs font-light', className)}
	{...restProps}
>
	{@render children?.()}
</Table.Head>
