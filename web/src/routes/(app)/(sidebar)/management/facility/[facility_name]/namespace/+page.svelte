<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	onMount(() => {
		const modelUuid = page.url.searchParams.get('scope') || '';
		const clusterName = page.params.facility_name;
		goto(`/management/facility/${clusterName}?scope=${modelUuid}`);
	});
</script>
