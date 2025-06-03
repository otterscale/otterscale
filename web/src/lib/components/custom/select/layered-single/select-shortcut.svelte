<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import type { HTMLAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();
</script>

<DropdownMenu.Shortcut bind:ref data-slot="select-shortcut" class={cn(className)} {...restProps}>
	{@render children?.()}
</DropdownMenu.Shortcut>
