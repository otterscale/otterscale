<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();
</script>

<DropdownMenu.Shortcut
	bind:ref
	data-slot="alert-dialog-action-shortcut"
	class={cn(className)}
	{...restProps}
>
	{@render children?.()}
</DropdownMenu.Shortcut>
