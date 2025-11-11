<script lang="ts">
	import BookmarkIcon from '@lucide/svelte/icons/bookmark';
	import type { Snippet } from 'svelte';

	import { page } from '$app/state';
	import { AppSidebar } from '$lib/components/layout';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Popover from '$lib/components/ui/popover';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import type { Path } from '$lib/path';
	import { bookmarks, breadcrumbs } from '$lib/stores';

	import type { LayoutData } from './$types';

	interface Props {
		data: LayoutData;
		children: Snippet;
	}

	let { data, children }: Props = $props();
	let open = $state(false);

	// Computed values
	const current = $derived($breadcrumbs.at(-1));
	const isBookmarked = $derived($bookmarks.some((bookmark) => bookmark.url === current?.url));

	// Event handlers
	function handleBookmarkAdd(path: Path | undefined) {
		if (!path) return;
		bookmarks.update((items) => [...items, path]);
		open = false;
	}

	function handleBookmarkDelete(path: Path | undefined) {
		if (!path) return;
		bookmarks.update((items) => items.filter((bookmark) => bookmark.url !== path.url));
		open = false;
	}
</script>

<svelte:head>
	<title>{current ? `${current.title} - ` : ''}OtterScale</title>
</svelte:head>

<Sidebar.Provider>
	<AppSidebar active={page.params.scope!} user={data.user} />
	<Sidebar.Inset>
		<header
			class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
		>
			<div class="flex w-full items-center justify-between gap-2 px-4">
				<!-- Sidebar Toggle -->
				<Sidebar.Trigger class="-ml-1 {buttonVariants({ variant: 'ghost', size: 'icon' })}" />
				<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />

				<!-- Breadcrumb Navigation -->
				<nav aria-label="Breadcrumb">
					<Breadcrumb.Root>
						<Breadcrumb.List>
							{#each $breadcrumbs as item}
								{#if item.url === current?.url}
									<Breadcrumb.Item>
										<Breadcrumb.Page>{current.title}</Breadcrumb.Page>
									</Breadcrumb.Item>
								{:else}
									<Breadcrumb.Item class="hidden md:block">
										<Breadcrumb.Link href={item.url}>
											{item.title}
										</Breadcrumb.Link>
									</Breadcrumb.Item>
									<Breadcrumb.Separator class="hidden md:block" />
								{/if}
							{/each}
						</Breadcrumb.List>
					</Breadcrumb.Root>
				</nav>

				<!-- Bookmark Popover -->
				<Popover.Root bind:open>
					<Popover.Trigger
						class="-mr-1 ml-auto {buttonVariants({ variant: 'ghost', size: 'icon' })}"
					>
						<BookmarkIcon fill={isBookmarked ? 'currentColor' : 'none'} />
						<span class="sr-only">Bookmark</span>
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
			</div>
		</header>

		<main class="flex flex-1 flex-col px-2 md:px-4 lg:px-8">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
