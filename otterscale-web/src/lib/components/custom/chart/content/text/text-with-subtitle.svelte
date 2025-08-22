<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

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
<div class="text-muted-foreground text-xs">
	{subtitle}
</div>
