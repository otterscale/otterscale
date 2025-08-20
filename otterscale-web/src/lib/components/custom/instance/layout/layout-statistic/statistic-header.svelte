<script lang="ts" module>
	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Header
	bind:ref
	data-slot="layout-statistic-header"
	class={cn('flex h-10 items-center', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Header>
