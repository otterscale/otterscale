<script lang="ts">
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import * as Card from '$lib/components/ui/card/index.js';

	let {
		title,
		description
	}: {
		title: string;
		description?: string;
	} = $props();
</script>

<Card.Root class="flex flex-col">
	<Card.Header class="flex flex-col items-center space-y-0 border-b p-0 sm:flex-row">
		<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
			<Title {title} />

			{#if description}
				<Card.Description>
					<Description {description} />
				</Card.Description>
			{/if}
		</div>
	</Card.Header>
	<Card.Content>
		<Content value="ERROR" />
	</Card.Content>
</Card.Root>
