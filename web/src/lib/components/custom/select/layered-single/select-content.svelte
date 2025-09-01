<script lang="ts" module>
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		sideOffset = 4,
		portalProps,
		class: className,
		align = 'center',
		side = 'bottom',
		...restProps
	}: DropdownMenuPrimitive.ContentProps & {
		portalProps?: DropdownMenuPrimitive.PortalProps;
	} = $props();
</script>

<DropdownMenu.Content
	{...portalProps}
	bind:ref
	data-slot="select-content"
	{sideOffset}
	{align}
	{side}
	class={cn('w-fit min-w-50 p-0', className)}
	{...restProps}
/>
