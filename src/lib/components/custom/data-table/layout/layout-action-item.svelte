<script lang="ts">
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		onSelect,
		...restProps
	}: DropdownMenuPrimitive.ItemProps & {
		inset?: boolean;
	} = $props();
</script>

<DropdownMenu.Item
	bind:ref
	class={cn(className)}
	onSelect={(e) => {
		e.preventDefault();
		onSelect?.(e);
	}}
	{...restProps}
/>
