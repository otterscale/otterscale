<script lang="ts">
	import { cn } from '$lib/utils';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		children,
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div bind:this={ref} class={cn('flex flex-col gap-4', className)} {...restProps}>
	{@render children?.()}
</div>
