<script lang="ts" module>
	import * as Tabs from '$lib/components/custom/tabs/naive/index';
	import { cn } from '$lib/utils.js';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import { IndexManager } from './utils.svelte';
</script>

<script lang="ts">
	let { ref = $bindable(null), class: className, ...restProps }: TabsPrimitive.ListProps = $props();

	setContext('IndexManager', new IndexManager());
</script>

<Tabs.List
	bind:ref
	data-slot="multiple-step-modal-steps"
	class={cn(
		'h-fit w-full items-end justify-around border-none transition-all duration-300',
		className
	)}
	{...restProps}
/>
