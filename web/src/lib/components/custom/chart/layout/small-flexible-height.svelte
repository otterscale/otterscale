<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { Snippet } from 'svelte';

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

<Card.Root class={className}>
	<Card.Header class="h-[30px] items-center">
		<Card.Title>
			{@render title()}
		</Card.Title>

		{#if description}
			<Card.Description>
				{@render description()}
			</Card.Description>
		{/if}

		{#if controller}
			{@render controller()}
		{/if}
	</Card.Header>

	<Card.Content>
		{@render content()}
	</Card.Content>

	<Card.Footer class="h-[0px]">
		{#if footer}
			{@render footer()}
		{/if}
	</Card.Footer>
</Card.Root>
