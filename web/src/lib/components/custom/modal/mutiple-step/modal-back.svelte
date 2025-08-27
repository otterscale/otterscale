<script lang="ts" module>
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { cn } from '$lib/utils.js';
	import { AlertDialog as AlertDialogPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
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

{#if !stepManager.isFirstStep}
	<AlertDialog.Action
		bind:ref
		data-slot="alert-dialog-action"
		class={cn(className)}
		{...restProps}
		onclick={(e) => {
			stepManager.back();
			onclick?.(e);
		}}
	/>
{/if}
