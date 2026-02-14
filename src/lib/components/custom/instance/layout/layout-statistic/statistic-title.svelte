<script lang="ts" module>
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
</script>

<script lang="ts">
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
