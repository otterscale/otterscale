<script lang="ts">
	import { onMount } from 'svelte';

	import { RecentContent, RecentSidebar } from '$lib/components/recent';
	import { listVisits, type pbVisit } from '$lib/pb';

	let items: pbVisit[] = [];
	onMount(async () => {
		items = await listVisits();
	});
</script>

<main
	class="flex min-h-[calc(100vh_-_theme(spacing.16))] flex-1 flex-col gap-4 bg-muted/40 p-4 md:gap-8 md:p-10"
>
	<div class="mx-auto grid w-full max-w-6xl gap-2">
		<div class="space-y-1">
			<h1 class="text-3xl font-semibold">Recently Viewed</h1>
			<p class="text-sm text-muted-foreground">
				View your browsing history and quickly access recently visited pages and resources.
			</p>
		</div>
	</div>
	<div
		class="mx-auto grid w-full max-w-6xl items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
	>
		<RecentSidebar {items} />
		<RecentContent {items} />
	</div>
</main>
