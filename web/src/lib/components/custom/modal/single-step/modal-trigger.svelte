<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { cn } from '$lib/utils';
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';
	import type { TriggerVariant } from './type';
	import { triggerVariants } from './utils';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		variant = 'default',
		class: className,
		...restProps
	}: AlertDialogPrimitive.TriggerProps & {
		variant?: TriggerVariant;
	} = $props();
</script>

<AlertDialog.Trigger
	bind:ref
	data-slot="modal-trigger"
	class={cn('capitalize', triggerVariants({ variant }), className)}
	{...restProps}
/>
