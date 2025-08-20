<script lang="ts">
	import { cn } from '$lib/utils';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		children,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div class={cn('text-muted-foreground', className)}>
	{@render children?.()}
</div>
