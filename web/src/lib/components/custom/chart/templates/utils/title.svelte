<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	let {
		title,
		class: className
	}: {
		title: Snippet;
		class?: string;
	} = $props();
</script>

<div class={cn('whitespace-nowrap', className)}>
	{@render title()}
</div>
