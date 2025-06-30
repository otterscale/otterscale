<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import type { Snippet } from 'svelte';
	import Hint from './utils/hint.svelte';
	import Title from './utils/title.svelte';

	let {
		title,
		content,
		description,
		footer,
		hint,
		class: className
	}: {
		title: Snippet;
		content: Snippet;
		description?: Snippet;
		footer?: Snippet;
		hint?: Snippet;
		class?: string;
	} = $props();
</script>

<Card.Root class="bg-muted/40 h-full w-full text-nowrap border-none shadow-none">
	<Card.Header class="h-[100px]">
		<Card.Title>
			<Title {title} class="text-3xl" />
			{#if hint}
				<Hint {hint} />
			{/if}
		</Card.Title>
		<Card.Description class="text-xl">
			{@render description?.()}
		</Card.Description>
	</Card.Header>
	<Card.Content class="h-[300px] text-7xl">
		{@render content()}
	</Card.Content>
	<Card.Footer class="text-xl">
		{@render footer?.()}
	</Card.Footer>
</Card.Root>
