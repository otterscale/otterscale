<script lang="ts">
	import type { TabsContentProps } from 'bits-ui';
	import * as Tabs from '$lib/components/custom/tabs/naive/index';
	import { cn } from '$lib/utils.js';
	import { getContext } from 'svelte';
	import type { IndexManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: Omit<TabsContentProps, 'value'> = $props();

	const indexManager: IndexManager = getContext('IndexManager');
</script>

<Tabs.Content bind:ref class={cn(className)} value={String(indexManager.get())} {...restProps} />
