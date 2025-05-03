<script lang="ts">
	import { onMount } from 'svelte';

	import { FavContent, FavSidebar } from '$lib/components/favorite';
	import { type favorite } from '$lib/models';

	let items: favorite[] = [];
	onMount(async () => {
		items = [
			{
				id: '1',
				path: '/dashboard',
				name: 'Dashboard Overview',
				created: new Date(Date.now() - 24 * 60 * 60 * 1000),
				updated: new Date()
			},
			{
				id: '2',
				path: '/settings',
				name: 'User Settings',
				created: new Date(Date.now() - 24 * 60 * 60 * 1000 * 2),
				updated: new Date()
			},
			{
				id: '3',
				name: 'Reports',
				path: '/reports',
				created: new Date(Date.now() - 24 * 60 * 60 * 1000 * 2),
				updated: new Date()
			}
		];
	});
</script>

<main
	class="flex min-h-[calc(100vh_-_theme(spacing.16))] flex-1 flex-col gap-4 bg-muted/40 p-4 md:gap-8 md:p-10"
>
	<div class="mx-auto grid w-full max-w-6xl gap-2">
		<div class="space-y-1">
			<h1 class="text-3xl font-semibold">Favorites</h1>
			<p class="text-sm text-muted-foreground">
				Browse and manage your pinned items, organized by features and functions.
			</p>
		</div>
	</div>
	<div
		class="mx-auto grid w-full max-w-6xl items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
	>
		<FavSidebar {items} />
		<FavContent {items} />
	</div>
</main>
