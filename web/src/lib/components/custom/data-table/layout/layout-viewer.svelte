<script lang="ts" module>
	import { cn } from '$lib/utils';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div
	bind:this={ref}
	class={cn('[&_th]:bg-muted rounded-lg border [&_th]:first:rounded-tl-lg [&_th]:last:rounded-tr-lg', className)}
	{...restProps}
>
	{@render children?.()}
</div>
