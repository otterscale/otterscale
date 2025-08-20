<script lang="ts">
	import type { WithElementRef } from 'bits-ui';
	import { getContext, setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { IndexManager, StepManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();

	const stepManager: StepManager = getContext('StepManager');
	setContext('IndexManager', new IndexManager(stepManager.steps));
</script>

<div bind:this={ref} class={className} {...restProps}>
	{@render children?.()}
</div>
