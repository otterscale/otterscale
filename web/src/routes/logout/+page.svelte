<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { logout } from '$lib/auth';
	import { Spinner } from '$lib/components/ui/spinner';
	import { m } from '$lib/paraglide/messages';

	onMount(async () => {
		await logout();
		await new Promise((resolve) => setTimeout(resolve, 1500));
		goto(resolve('/'));
	});
</script>

<div class="h-svh">
	<div class="m-auto flex h-full w-1/4 flex-col items-center justify-center gap-2">
		<h1 class="text-9xl leading-tight font-bold">302</h1>
		<p class="text-muted-foreground line-clamp-2 text-center">
			{m.redirect_message_302()}
		</p>
		<div class="mt-6 flex gap-4">
			<Spinner class="size-8" />
		</div>
	</div>
</div>
