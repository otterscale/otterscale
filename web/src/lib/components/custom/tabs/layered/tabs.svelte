<script lang="ts" module>
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';

	import * as Tabs from '$lib/components/ui/tabs';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
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
