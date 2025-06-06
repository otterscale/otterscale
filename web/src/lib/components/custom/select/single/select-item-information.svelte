<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import Icon from '@iconify/svelte';

	import { Tooltip as TooltipPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		sideOffset = 4,
		children,
		...restProps
	}: TooltipPrimitive.ContentProps & {} = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger>
			<Icon icon="ph:info" />
		</Tooltip.Trigger>
		<Tooltip.Content
			bind:ref
			{sideOffset}
			class={cn('bg-popover text-popover-foreground border p-2 shadow', className)}
			{...restProps}
		>
			{@render children?.()}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
