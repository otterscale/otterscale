<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
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
