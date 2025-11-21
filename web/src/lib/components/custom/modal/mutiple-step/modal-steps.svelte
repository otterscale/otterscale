<script lang="ts" module>
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { getContext, setContext } from 'svelte';

	import * as Tabs from '$lib/components/ui/tabs';
	import { cn } from '$lib/utils.js';

	import { IndexManager, StepManager } from './utils.svelte';
</script>

<script lang="ts">
	let { ref = $bindable(null), class: className, ...restProps }: TabsPrimitive.ListProps = $props();

	const stepManager: StepManager = getContext('StepManager');
	stepManager.reset();
	setContext('IndexManager', new IndexManager(stepManager.steps));
</script>

<Tabs.List
	bind:ref
	data-slot="multiple-step-modal-steps"
	class={cn(
		'h-fit w-full items-start justify-around border-none p-4 transition-all duration-300',
		className
	)}
	{...restProps}
/>
