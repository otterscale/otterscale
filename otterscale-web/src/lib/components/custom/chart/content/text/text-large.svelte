<!-- <script lang="ts">
	import { cn } from '$lib/utils';
	import * as Card from '$lib/components/ui/card';

	let {
		value,
		unit,
		class: className
	}: {
		value?: string | number;
		unit?: string;
		class?: string;
	} = $props();
</script>

<Card.Content class={cn('text-7xl', className)}>
	{value || 'No Data'}
	{unit || ''}
</Card.Content> -->

<script lang="ts">
	import { cn } from '$lib/utils';
	import * as Card from '$lib/components/ui/card';

	let {
		value,
		unit,
		children,
		class: className
	}: {
		value?: string | number;
		unit?: string;
		children?: import('svelte').Snippet;
		class?: string;
	} = $props();
</script>

<Card.Content class={cn('flex items-center gap-2 text-7xl', className)}>
	{#if children}
		{@render children()}
	{:else}
		{value !== undefined && value !== null ? value : 'No Data'}
	{/if}
	{unit || ''}
</Card.Content>
