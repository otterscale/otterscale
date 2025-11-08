<script lang="ts" module>
	import { type AlertDialogRootPropsWithoutHTML } from 'bits-ui';
	import { setContext } from 'svelte';

	import type { StepManagerState } from './types';
	import { StepManager } from './utils.svelte';

	import * as AlertDialog from '$lib/components/custom/alert-dialog';
</script>

<script lang="ts">
	let {
		open = $bindable(false),
		steps,
		...restProps
	}: AlertDialogRootPropsWithoutHTML & { steps: number } = $props();

	let isUpdating = false;
	setContext('Accessor', {
		set open(value: boolean) {
			open = value;
		}
	});
	setContext(
		'StepManager',
		new StepManager(steps, {
			get isUpdating() {
				return isUpdating;
			},
			set isUpdating(value: boolean) {
				isUpdating = value;
			}
		} as StepManagerState)
	);
</script>

<AlertDialog.Root bind:open {...restProps} />
