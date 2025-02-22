<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { SiteFooter, SiteHeader } from '$lib/components';
	import pb from '$lib/pb';

	onMount(() => {
		if (page.url.pathname == '/') {
			return;
		}
		if (page.url.pathname == '/logout') {
			return;
		}
		if (!pb.authStore.isValid) {
			goto('/login?callback=' + page.url.pathname);
		}
	});
</script>

<div class="relative flex min-h-screen flex-col bg-background" data-vaul-drawer-wrapper>
	<SiteHeader />
	<div class="flex-1">
		<slot />
	</div>
	<SiteFooter />
</div>
