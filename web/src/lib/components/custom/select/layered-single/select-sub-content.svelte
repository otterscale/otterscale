<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		align = 'start',
		side = 'right',
		...restProps
	}: DropdownMenuPrimitive.SubContentProps = $props();
</script>

<DropdownMenu.SubContent
	bind:ref
	data-slot="select-sub-content"
	class={cn('p-0', className)}
	{align}
	{side}
	{...restProps}
/>
