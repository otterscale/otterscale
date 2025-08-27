<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import Title from '$lib/components/custom/chart/title.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';

	let {
		title,
		description,
		class: className
	}: {
		title: string;
		description?: string;
		class?: string;
	} = $props();
</script>

<Card.Root class="flex flex-col">
	<Card.Header class="items-center">
		<Title {title} />

		{#if description}
			<Card.Description>
				<Description {description} />
			</Card.Description>
		{/if}
	</Card.Header>
	<Card.Content>
		<Content value="ERROR" />
	</Card.Content>
</Card.Root>
