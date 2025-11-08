<script lang="ts" module>
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	import { type ButtonVariant,buttonVariants } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		sideOffset = 4,
		class: className,
		label: label = 'Actions',
		variant: variant = 'default',
		...restProps
	}: DropdownMenuPrimitive.ContentProps & { label?: string; variant?: ButtonVariant } = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger class={cn('hover:cursor-pointer', buttonVariants({ variant: variant }))}>
		{label}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content
		bind:ref
		data-slot="alert-dialog-actions"
		side="right"
		align="start"
		{sideOffset}
		class={cn(className)}
		{...restProps}
	/>
</DropdownMenu.Root>
