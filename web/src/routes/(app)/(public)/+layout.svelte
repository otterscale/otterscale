<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import pb from '$lib/pb';

	onMount(() => {
		if (pb.authStore.isValid) {
			goto('/');
		}
	});
</script>

<slot />
