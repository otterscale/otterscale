<script lang="ts">
	import * as Command from '$lib/components/ui/command';

	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();
</script>

<Command.Shortcut bind:ref data-slot="select-shortcut" class={cn(className)} {...restProps}>
	{@render children?.()}
</Command.Shortcut>
