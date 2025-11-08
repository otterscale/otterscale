<script lang="ts" module>
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { cn } from '$lib/utils.js';

	import type { StepManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		onclick,
		...restProps
	}: AlertDialogPrimitive.ActionProps = $props();

	const stepManager: StepManager = getContext('StepManager');
</script>

{#if !stepManager.isLastStep}
	<AlertDialog.Action
		bind:ref
		data-slot="multiple-step-modal-next"
		class={cn(className)}
		{...restProps}
		onclick={(e) => {
			stepManager.next();
			onclick?.(e);
		}}
	/>
{/if}
