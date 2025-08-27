<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
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
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();
</script>

<Command.Shortcut bind:ref data-slot="select-shortcut" class={cn(className)} {...restProps}>
	{@render children?.()}
</Command.Shortcut>
