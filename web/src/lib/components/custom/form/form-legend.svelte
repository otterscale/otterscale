<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLLegendElement>> = $props();
</script>

<legend bind:this={ref} data-slot="form-legend" class={cn('text-lg font-bold', className)} {...restProps}>
	{@render children?.()}
</legend>
