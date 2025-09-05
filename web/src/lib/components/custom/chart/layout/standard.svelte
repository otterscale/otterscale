<script lang="ts">
	import type { Snippet } from 'svelte';

	import * as Card from '$lib/components/ui/card/index.js';
	import { cn } from '$lib/utils';

	let {
		title,
		description,
		controller,
		content,
		footer,
		class: className,
	}: {
		title: Snippet;
		description?: Snippet;
		controller?: Snippet;
		content: Snippet;
		footer?: Snippet;
		class?: string;
	} = $props();
</script>

<Card.Root class={cn('flex flex-col', className)}>
	<Card.Header class="flex flex-col items-center space-y-0 border-b p-0 sm:flex-row ">
		<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
			<!-- <div class="flex flex-1 flex-col justify-center px-6"> -->
			<Card.Title>
				{@render title()}
			</Card.Title>
			<Card.Description>
				{#if description}
					{@render description()}
				{:else}
					<span class="invisible">No description</span>
				{/if}
			</Card.Description>
		</div>
		{#if controller}
			{@render controller()}
		{/if}
	</Card.Header>
	<Card.Content>
		{@render content()}
	</Card.Content>
	{#if footer}
		<Card.Footer>
			{@render footer()}
		</Card.Footer>
	{/if}
</Card.Root>
