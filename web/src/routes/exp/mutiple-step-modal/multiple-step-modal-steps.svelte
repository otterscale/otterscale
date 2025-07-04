<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs/naive/index';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { setContext } from 'svelte';
	import { IndexManager } from './utils.svelte';

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
