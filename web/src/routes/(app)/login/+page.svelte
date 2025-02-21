<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import pb from '$lib/pb';

	onMount(() => {
		if (pb.authStore.isValid) {
			goto('/');
		}
	});

	function callback() {
		const callbackParam = page.url.searchParams.get('callback');
		if (callbackParam) {
			goto(callbackParam);
			return;
		}
		goto('/');
	}

	async function github() {
		await pb.collection('users').authWithOAuth2({ provider: 'github' });
		callback();
	}
</script>

<Button on:click={github}>
	<Icon icon="ph:github-logo" class="h-5 w-5" />
	Github
</Button>
