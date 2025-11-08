<script lang="ts" module>
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
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
