<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import pb from '$lib/pb';

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

<Button on:click={github}>Github</Button>
