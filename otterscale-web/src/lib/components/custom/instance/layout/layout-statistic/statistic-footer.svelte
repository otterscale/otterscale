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

<Card.Footer
	bind:ref
	data-slot="layout-statistic-footer"
	class={cn('text-muted-foreground h-10 items-end truncate text-xs font-light', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Footer>
