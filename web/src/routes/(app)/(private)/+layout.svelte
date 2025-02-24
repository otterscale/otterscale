<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteFooter, SiteHeader } from '$lib/components';
	import pb, { upsertVisit } from '$lib/pb';
	import { setCallback } from '$lib/utils';
	import { i18n } from '$lib/i18n';

	let { children } = $props();

	let currentPage = i18n.route(page.url.pathname);

	onMount(async () => {
		await upsertVisit();
	});

	$effect(() => {
		if (currentPage !== i18n.route(page.url.pathname)) {
			(async () => await upsertVisit())();
			currentPage = i18n.route(page.url.pathname);
		}
		if (i18n.route(page.url.pathname) == '/' || i18n.route(page.url.pathname) == '/logout') {
			return;
		}
		if (!pb.authStore.isValid) {
			goto(setCallback(i18n.resolveRoute('/login')));
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
