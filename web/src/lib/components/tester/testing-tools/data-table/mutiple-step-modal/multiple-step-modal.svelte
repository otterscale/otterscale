<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog/index';
	import { type AlertDialogRootPropsWithoutHTML } from 'bits-ui';
	import { setContext } from 'svelte';
	import { StepManager } from './utils.svelte';

	let {
		open = $bindable(false),
		steps,
		...restProps
	}: AlertDialogRootPropsWithoutHTML & { steps: number } = $props();

	setContext('Accessor', {
		set open(value: boolean) {
			open = value;
		}
	});
	setContext('StepManager', new StepManager(3));
</script>

<AlertDialog.Root bind:open {...restProps} />
