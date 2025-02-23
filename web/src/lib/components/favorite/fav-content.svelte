<script lang="ts">
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import type { pbFavorite } from '$lib/pb';

	export let items: pbFavorite[];

	function filter(): pbFavorite[] {
		if (page.url.hash === '') {
			return items;
		}
		return items.filter((item) => {
			var paths = item.path.split('/');
			return page.url.hash == `#${paths[1]}`;
		});
	}
</script>

<div class="grid gap-6">
	{#each filter() as item}
		<Card.Root>
			<Card.Content class="space-y-1 py-4">
				<Card.Title>{item.name}</Card.Title>
				<Card.Description>{item.path}</Card.Description>
			</Card.Content>
		</Card.Root>
	{/each}
</div>
