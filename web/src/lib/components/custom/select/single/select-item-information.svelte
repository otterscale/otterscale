<script lang="ts" module>
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Tooltip as TooltipPrimitive } from 'bits-ui';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		sideOffset = 4,
		children,
		...restProps
	}: TooltipPrimitive.ContentProps & {} = $props();
</script>

<div data-slot="select-item-information">
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
</div>
