<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteFooter, SiteHeader } from '$lib/components';
	import pb, { upsertVisit } from '$lib/pb';
	import { setCallback } from '$lib/utils';

	let { children } = $props();

	onMount(async () => {
		await upsertVisit();
	});

	let currentPage = page.url.pathname;

	$effect(() => {
		if (currentPage !== page.url.pathname) {
			(async () => await upsertVisit())();
		}
		currentPage = page.url.pathname;
		if (currentPage == '/' || currentPage == '/logout') {
			return;
		}
		if (!pb.authStore.isValid) {
			goto(setCallback('/login'));
		}
	});
</script>

<div class="relative flex min-h-screen flex-col bg-background" data-vaul-drawer-wrapper>
	<SiteHeader />
	<div class="flex-1">
		{@render children()}
	</div>
	<SiteFooter />
</div>
