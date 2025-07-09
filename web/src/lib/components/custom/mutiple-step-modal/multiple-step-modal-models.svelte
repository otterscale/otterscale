<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
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
