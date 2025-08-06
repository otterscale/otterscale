<script lang="ts">
	import { applicationsPath, applicationsStorePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import Store from '$lib/components/applications/store/index.svelte';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [applicationsPath], current: applicationsStorePath });
</script>

{#if $activeScope}
	<Store />
{/if}
