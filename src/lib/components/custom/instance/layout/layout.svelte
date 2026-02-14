<script lang="ts" module>
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div class={cn('space-y-4', className)}>
	{@render children?.()}
</div>
