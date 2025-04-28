<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	onMount(() => {
		const modelUuid = $page.params.scope_uuid;
		goto(`/management/scope/${modelUuid}`);
	});
</script>
