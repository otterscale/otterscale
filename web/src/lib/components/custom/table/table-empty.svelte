<script lang="ts" module>
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className,
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

<div class={cn('text-muted-foreground p-2 text-center', className)}>
	{#if children}
		{@render children()}
	{:else}
		No Data
	{/if}
</div>
