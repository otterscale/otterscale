<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { cn } from '$lib/utils.js';
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { StepManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		onclick,
		...restProps
	}: AlertDialogPrimitive.ActionProps = $props();

	const stepManager: StepManager = getContext('StepManager');
</script>

<AlertDialog.Cancel
	bind:ref
	data-slot="alert-dialog-cancel"
	class={cn(className)}
	{...restProps}
	onclick={(e) => {
		stepManager.reset();
		onclick?.(e);
	}}
/>
