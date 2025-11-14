<script lang="ts">
	import BookmarkIcon from '@lucide/svelte/icons/bookmark';

	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Popover from '$lib/components/ui/popover';
	import { m } from '$lib/paraglide/messages';
	import type { Path } from '$lib/path';
	import { bookmarks, breadcrumbs } from '$lib/stores';

	const current = $derived($breadcrumbs.at(-1));
	const isBookmarked = $derived($bookmarks.some((bookmark) => bookmark.url === current?.url));

	let open = $state(false);
	function close() {
		open = false;
	}

	function handleBookmarkAdd(path: Path | undefined) {
		if (!path) return;
		bookmarks.update((items) => [...items, path]);
		close();
	}

	function handleBookmarkDelete(path: Path | undefined) {
		if (!path) return;
		bookmarks.update((items) => items.filter((bookmark) => bookmark.url !== path.url));
		close();
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger class="-mr-1 ml-auto {buttonVariants({ variant: 'ghost', size: 'icon' })}">
		<BookmarkIcon fill={isBookmarked ? 'currentColor' : 'none'} />
		<span class="sr-only">{m.bookmark()}</span>
	</Popover.Trigger>

	<Popover.Content align="start" side="left">
		<div class="grid gap-4 p-2">
			<div class="space-y-1">
				<h4 class="leading-none font-medium">{m.bookmark_added()}</h4>
			</div>

			<div class="grid gap-2">
				<div class="grid grid-cols-3 items-center gap-4">
					<Label for="bookmark-name">{m.name()}</Label>
					<Input id="bookmark-name" value={current?.title} class="col-span-2 h-8" />
				</div>
			</div>

			<div class="grid grid-cols-2 gap-6">
				<Button variant="secondary" onclick={() => handleBookmarkDelete(current)}>
					{m.remove()}
				</Button>
				<Button onclick={() => handleBookmarkAdd(current)}>
					{m.complete()}
				</Button>
			</div>
		</div>
	</Popover.Content>
</Popover.Root>
