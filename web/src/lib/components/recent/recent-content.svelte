<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import { i18n } from '$lib/i18n';
	import type { recent } from '$lib/models';
	import { formatTimeAgo } from '$lib/formatter';

	export let items: recent[];

	function filter(): recent[] {
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
			class="inline-flex whitespace-nowrap rounded-md shadow transition-colors hover:cursor-pointer hover:text-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
			onclick={() => goto(i18n.resolveRoute(item.path))}
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
