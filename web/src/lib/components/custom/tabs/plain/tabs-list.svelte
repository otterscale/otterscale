<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let { ref = $bindable(null), class: className, ...restProps }: TabsPrimitive.ListProps = $props();
</script>

<Tabs.List
	bind:ref
	data-slot="tabs-list"
	class={cn('flex w-full items-end justify-start rounded-none border-b bg-transparent', className)}
	{...restProps}
/>
