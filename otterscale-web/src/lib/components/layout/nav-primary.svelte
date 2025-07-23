<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import Icon from '@iconify/svelte';
	import { getIconFromUrl } from './icon';

	interface Bookmark {
		name: string;
		url: string;
	}

	let { bookmarks }: { bookmarks: Bookmark[] } = $props();

	let visibleCount = $state(3);
	const increment = 3;

	const visibleBookmarks = $derived(bookmarks.slice(0, visibleCount));
	const hasMoreBookmarks = $derived(bookmarks.length > visibleCount);

	function deleteBookmark(bookmark: Bookmark): void {
		toast.warning(`TODO: delete ${bookmark.name.toLowerCase()}`);
	}

	function showMoreBookmarks(): void {
		visibleCount += increment;
	}
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>Bookmarks</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each visibleBookmarks as bookmark (bookmark.name)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<a href={bookmark.url} {...props}>
							<Icon icon={getIconFromUrl(bookmark.url)} />
							<span>{bookmark.name}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
				<Sidebar.MenuAction showOnHover onclick={() => deleteBookmark(bookmark)}>
					<Icon icon="ph:x-bold" class="text-red-500" />
					<span class="sr-only">Delete {bookmark.name}</span>
				</Sidebar.MenuAction>
			</Sidebar.MenuItem>
		{/each}

		{#if hasMoreBookmarks}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton onclick={showMoreBookmarks}>
					<Icon icon="ph:dots-three-bold" />
					More
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}
	</Sidebar.Menu>
</Sidebar.Group>
