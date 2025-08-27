<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
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

<div
	bind:this={ref}
	data-slot="form-description"
	class={cn('text-muted-foreground -mt-2 select-none text-sm', className)}
	{...restProps}
>
	{@render children?.()}
</div>
