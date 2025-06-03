<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: DropdownMenuPrimitive.SeparatorProps = $props();
</script>

<DropdownMenu.Separator
	bind:ref
	data-slot="select-separator"
	class={cn(className)}
	{...restProps}
/>
