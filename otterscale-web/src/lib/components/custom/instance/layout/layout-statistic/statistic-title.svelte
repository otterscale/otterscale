<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Title
	bind:ref
	data-slot="layout-statistic-title"
	class={cn('flex items-center gap-1', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Title>
