<script lang="ts">
	import { Tabs as TabsPrimitive } from 'bits-ui';

	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { cn } from '$lib/utils.js';

	let { ref = $bindable(null), class: className, ...restProps }: TabsPrimitive.TriggerProps = $props();
</script>

<Tabs.Trigger
	bind:ref
	data-slot="tabs-trigger"
	class={cn(
		'data-[state=active]:border-b-muted-foreground data-[state=active]:text-primary text-muted-foreground -mb-[4px] w-fit rounded-none hover:cursor-pointer data-[state=active]:border-b data-[state=active]:shadow-none',
		className,
	)}
	{...restProps}
/>
