<script lang="ts">
	import { page } from '$app/state';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: 'NTP Server', url: '' }
	});
</script>

<Card.Root class="shadow-none">
	<Card.Content>
		<h3 class="text-2xl tracking-tight">This page is currently under construction.</h3>
		<p class="text-muted-foreground text-sm">Please check back later.</p>
		<Button class="mt-4" href="/">Home</Button>
	</Card.Content>
</Card.Root>
