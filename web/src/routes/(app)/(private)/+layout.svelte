<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteFooter, SiteHeader } from '$lib/components';
	import pb, { upsertVisit } from '$lib/pb';
	import { setCallback } from '$lib/utils';

	let { children } = $props();

	let currentPage = page.url.pathname;

	onMount(async () => {
		await upsertVisit();
	});

	$effect(() => {
		if (currentPage !== page.url.pathname) {
			(async () => await upsertVisit())();
			currentPage = page.url.pathname;
		}
		if (page.url.pathname == '/' || page.url.pathname == '/logout') {
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
