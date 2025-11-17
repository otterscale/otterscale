<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { getPathIcon, type Path } from '$lib/path';
	import { bookmarks } from '$lib/stores';

	const increment = 3;
	let visibleCount = $state(3);

	const visibleBookmarks = $derived($bookmarks.slice(0, visibleCount));
	const hasMoreBookmarks = $derived($bookmarks.length > visibleCount);

	function showMoreBookmarks(): void {
		visibleCount += increment;
	}

	async function onBookmarkDelete(path: Path) {
		bookmarks.update((currentBookmarks) =>
			currentBookmarks.filter((bookmark) => bookmark.url !== path.url)
		);
	}
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>{m.bookmarks()}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each visibleBookmarks as bookmark (bookmark.title)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<!-- eslint-disable svelte/no-navigation-without-resolve -->
						<a href={bookmark.url} {...props}>
							<Icon icon={getPathIcon(bookmark.url)} />
							<span>{bookmark.title}</span>
						</a>
						<!-- eslint-enable svelte/no-navigation-without-resolve -->
					{/snippet}
				</Sidebar.MenuButton>
				<Sidebar.MenuAction showOnHover onclick={() => onBookmarkDelete(bookmark)}>
					<Icon icon="ph:x-bold" class="text-red-500" />
					<span class="sr-only">Delete {bookmark.title}</span>
				</Sidebar.MenuAction>
			</Sidebar.MenuItem>
		{/each}

		{#if hasMoreBookmarks}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton onclick={showMoreBookmarks}>
					<Icon icon="ph:dots-three-bold" />
					{m.more()}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}

		{#if $bookmarks.length === 0}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton aria-disabled>
					<Icon icon="ph:empty" />
					{m.empty_bookmark()}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}
	</Sidebar.Menu>
</Sidebar.Group>
