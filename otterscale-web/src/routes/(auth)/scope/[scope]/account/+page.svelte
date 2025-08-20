<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.account(page.params.scope) });
</script>

<div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center gap-2">
	<div class="border-border flex size-28 items-center justify-center rounded-full border-2">
		<Icon icon="ph:rocket-launch" class="size-16" />
	</div>

	<p class="pt-2 text-3xl">{m.under_development()}</p>
	<p class="text-muted-foreground text-base">
		{m.under_development_description()}
	</p>
</div>
