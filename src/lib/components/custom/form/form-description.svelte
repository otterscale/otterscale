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
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	data-slot="form-description"
	class={cn('-mt-2 text-sm text-muted-foreground select-none', className)}
	{...restProps}
>
	{@render children?.()}
</div>
