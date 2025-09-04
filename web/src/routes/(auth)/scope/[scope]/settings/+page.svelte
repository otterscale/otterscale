<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';

	onMount(() => {
		// Redirect on component mount
		goto(`${dynamicPaths.settings(page.params.scope).url}/ntp-server`);
	});
</script>
