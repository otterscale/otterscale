<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { setContext } from 'svelte';

	let {
		ref = $bindable(null),
		class: className,
		value = $bindable(),
		...restProps
	}: TabsPrimitive.RootProps = $props();

	const setter = (newValue: string | undefined) => {
		value = newValue;
	};
	setContext('setter', setter);
</script>

<Tabs.Root bind:ref class={cn(className)} {value} {...restProps} />
