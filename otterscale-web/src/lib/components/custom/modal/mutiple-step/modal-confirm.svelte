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
	const accessor: { open: boolean } = getContext('Accessor');
</script>

{#if stepManager.isLastStep}
	<AlertDialog.Action
		bind:ref
		data-slot="multiple-step-modal-confirm"
		class={cn(className)}
		{...restProps}
		onclick={(e) => {
			accessor.open = false;
			stepManager.reset();
			onclick?.(e);
		}}
	/>
{/if}
