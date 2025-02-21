<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import pb from '$lib/pb';

	onMount(() => {
		if (!pb.authStore.isValid) {
			goto('/account/login?callback=' + page.url.pathname);
		}
	});
</script>

<slot />
