<script lang="ts">
	import type { Snippet } from 'svelte';

	import { cn } from '$lib/utils';

	let {
		value,
		unit,
		subtitle,
		children,
		class: className
	}: {
		value?: string | number;
		unit?: string;
		subtitle?: string;
		children?: Snippet;
		class?: string;
	} = $props();
</script>

<div class={cn('text-3xl', className)}>
	{#if children}
		{@render children()}
	{:else}
		{value !== undefined && value !== null ? value : 'No Data'}
	{/if}
	{unit || ''}
</div>
<div class="text-xs text-muted-foreground">
	{subtitle}
</div>
