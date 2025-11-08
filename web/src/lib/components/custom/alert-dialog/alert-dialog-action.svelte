<script lang="ts" module>
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import {
		type ButtonSize,
		type ButtonVariant,
		buttonVariants
	} from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
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
		'capitalize hover:cursor-pointer',
		buttonVariants({ variant: actionVariant, size: actionSize }),
		className
	)}
	{...restProps}
/>
