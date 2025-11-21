<script lang="ts">
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import * as Card from '$lib/components/ui/card';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Header
	bind:ref
	data-slot="data-table-statistics-header"
	class={cn('gap-4', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Header>
