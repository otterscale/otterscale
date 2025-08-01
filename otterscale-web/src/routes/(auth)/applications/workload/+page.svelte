<script lang="ts">
	import { applicationsPath, applicationsWorkloadPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [applicationsPath], current: applicationsWorkloadPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
