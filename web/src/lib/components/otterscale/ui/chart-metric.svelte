<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';

	export let description: string;
	export let metric: string | number;
	export let footer: string | undefined = undefined;
	export let tag: string | undefined = undefined;
</script>

<Card.Root>
	<Card.Header>
		<Card.Title
			><div class="flex justify-between text-xs">
				{description}
				<div class="font-light text-muted-foreground">
					{tag}
				</div>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="text-2xl">
			{metric}
		</div>
	</Card.Content>
	<Card.Footer>
		<div class="truncate text-xs font-light text-muted-foreground">
			{footer}
		</div>
	</Card.Footer>
</Card.Root>
