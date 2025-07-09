<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { IndexManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();

	setContext('IndexManager', new IndexManager());
</script>

<div bind:this={ref} class={className} {...restProps}>
	{@render children?.()}
</div>
