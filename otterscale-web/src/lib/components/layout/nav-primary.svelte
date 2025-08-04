<script lang="ts">
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { urlIcon } from '$lib/path';

	interface Bookmark {
		title: string;
		url: string;
	}

	let { bookmarks }: { bookmarks: Bookmark[] } = $props();

	let visibleCount = $state(3);
	const increment = 3;

	const visibleBookmarks = $derived(bookmarks.slice(0, visibleCount));
	const hasMoreBookmarks = $derived(bookmarks.length > visibleCount);

	function deleteBookmark(bookmark: Bookmark): void {
		toast.warning(`TODO: delete ${bookmark.title.toLowerCase()}`);
	}

	function showMoreBookmarks(): void {
		visibleCount += increment;
	}
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>{m.bookmarks()}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each visibleBookmarks as bookmark (bookmark.title)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<a href={bookmark.url} {...props}>
							<Icon icon={urlIcon(page.params.scope, bookmark.url)} />
							<span>{bookmark.title}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
				<Sidebar.MenuAction showOnHover onclick={() => deleteBookmark(bookmark)}>
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
	</Sidebar.Menu>
</Sidebar.Group>
