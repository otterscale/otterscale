<script lang="ts" module>
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';

	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: AlertDialogPrimitive.CancelProps = $props();
</script>

<AlertDialog.Cancel
	bind:ref
	data-slot="alert-dialog-cancel"
	class={cn('capitalize hover:cursor-pointer', buttonVariants({ variant: 'outline' }), className)}
	{...restProps}
/>
