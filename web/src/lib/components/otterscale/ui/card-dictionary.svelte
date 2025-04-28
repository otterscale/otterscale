<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';

	export let title: string | undefined = undefined;
</script>

<Card.Root>
	{#if title}
		<Card.Header>
			<Card.Title><div class="text-base">{title}</div></Card.Title>
		</Card.Header>
	{/if}
	<Card.Content>
		<slot />
	</Card.Content>
</Card.Root>
