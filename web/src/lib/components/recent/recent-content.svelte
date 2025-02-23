<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import type { pbVisit } from '$lib/pb';
	import { formatTimeAgo } from '$lib/utils';
	import Button from '../ui/button/button.svelte';

	export let items: pbVisit[];

	function filter(): pbVisit[] {
		if (page.url.hash === '') {
			return items;
		}
		return items.filter((item) => {
			var paths = item.path.split('/');
			return page.url.hash == `#${paths[1]}`;
		});
	}
</script>

<div class="grid gap-3">
	{#each filter() as item}
		<Card.Root
			class="inline-flex whitespace-nowrap rounded-md shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
			on:click={() => goto(item.path)}
		>
			<Card.Header class="pb-6">
				<Card.Title>
					{item.name}
				</Card.Title>
				<Card.Description class="text-xs">
					{item.path}
					·
					{formatTimeAgo(item.visited)}
				</Card.Description>
			</Card.Header>
		</Card.Root>
	{/each}
</div>
