<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		inset,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		inset?: boolean;
	} = $props();
</script>

<DropdownMenu.Label
	bind:ref
	data-slot="select-label"
	data-inset={inset}
	class={cn(className)}
	{...restProps}
>
	{@render children?.()}
</DropdownMenu.Label>
