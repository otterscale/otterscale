<script lang="ts" module>
	import { Popover as PopoverPrimitive } from 'bits-ui';

	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		sideOffset = 4,
		align = 'center',
		side = 'bottom',
		...restProps
	}: PopoverPrimitive.ContentProps & {
		portalProps?: PopoverPrimitive.PortalProps;
	} = $props();
</script>

<Popover.Content
	bind:ref
	data-slot="select-content"
	{sideOffset}
	{align}
	{side}
	class={cn('w-fit min-w-50 p-0', className)}
	{...restProps}
/>
