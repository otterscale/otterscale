<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';

	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';
	import {
		buttonVariants,
		type ButtonVariant,
		type ButtonSize
	} from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		actionVariant,
		actionSize,
		class: className,
		...restProps
	}: AlertDialogPrimitive.ActionProps & {
		actionVariant?: ButtonVariant;
		actionSize?: ButtonSize;
	} = $props();
</script>

<AlertDialog.Action
	bind:ref
	data-slot="alert-dialog-action"
	class={cn(
		'hover:cursor-pointer',
		buttonVariants({ variant: actionVariant, size: actionSize }),
		className
	)}
	{...restProps}
/>
