<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	onMount(() => {
		const modelUuid = $page.params.scope_uuid;
		const clusterName = $page.params.facility_name;
		goto(`/management/scope/${modelUuid}/facility/${clusterName}`);
	});
</script>
