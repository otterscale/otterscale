<script lang="ts">
	import { page } from '$app/state';
	import type { pbRecent } from '$lib/pb';

	export let items: pbRecent[];

	function features(): string[] {
		return [
			...new Set(
				items.map((item) => {
					var paths = item.path.split('/');
					return paths[paths.length - 1];
				})
			)
		]
			.filter(Boolean)
			.sort();
	}

	function active(v: string): string {
		return (v == '#' && page.url.hash === '') || page.url.hash === `#${v}`
			? 'font-semibold text-primary'
			: '';
	}
</script>

<nav
	class="grid gap-3 text-sm text-muted-foreground"
	data-x-chunk-container="chunk-container after:right-0"
>
	<a href="/recents" class={active('#')}>All</a>
	{#each features() as feature}
		<a href="#{feature}" class={active(feature)}>
			{feature.charAt(0).toUpperCase() + feature.slice(1)}
		</a>
	{/each}
</nav>
