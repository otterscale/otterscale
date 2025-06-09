<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import type { HTMLAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

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
