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
	data-slot="form-actions"
	class={cn('flex items-center justify-end gap-2', className)}
	{...restProps}
>
	{@render children?.()}
</div>
