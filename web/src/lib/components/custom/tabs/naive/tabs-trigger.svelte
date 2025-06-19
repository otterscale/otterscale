<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: TabsPrimitive.TriggerProps = $props();
</script>

<Tabs.Trigger
	bind:ref
	data-slot="tabs-trigger"
	class={cn(
		'data-[state=active]:border-b-muted-foreground -mb-[6px] w-fit rounded-none hover:cursor-pointer data-[state=active]:border-b data-[state=active]:shadow-none',
		className
	)}
	{...restProps}
/>
